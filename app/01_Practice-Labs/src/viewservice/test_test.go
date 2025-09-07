package viewservice

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"
)

// check validates that the current view matches the expected primary, backup, and view number.
// It performs multiple assertions to ensure the view service state is correct:
//   - Primary server matches expected value
//   - Backup server matches expected value  
//   - View number matches expected value (if n != 0)
//   - Clerk's Primary() method returns the correct primary
func check(t *testing.T, ck *Clerk, expectedPrimary string, expectedBackup string, expectedViewNum uint) {
	view, ok := ck.Get()
	if !ok {
		t.Fatalf("failed to get view from view service")
	}
	
	if view.Primary != expectedPrimary {
		t.Fatalf("wanted primary %v, got %v", expectedPrimary, view.Primary)
	}
	if view.Backup != expectedBackup {
		t.Fatalf("wanted backup %v, got %v", expectedBackup, view.Backup)
	}
	if expectedViewNum != 0 && expectedViewNum != view.Viewnum {
		t.Fatalf("wanted viewnum %v, got %v", expectedViewNum, view.Viewnum)
	}
	if ck.Primary() != expectedPrimary {
		t.Fatalf("wanted primary %v, got %v", expectedPrimary, ck.Primary())
	}
}

// port generates a unique Unix socket path for testing.
// It creates a path under /var/tmp/824-{uid}/viewserver-{pid}-{suffix}
// to ensure each test run uses unique socket names and avoids conflicts.
func port(suffix string) string {
	uid := strconv.Itoa(os.Getuid())
	pid := strconv.Itoa(os.Getpid())
	
	dir := "/var/tmp/824-" + uid + "/"
	os.Mkdir(dir, 0777)
	
	return dir + "viewserver-" + pid + "-" + suffix
}

// Test1 is the main test function that validates the view service functionality.
// It tests various scenarios including:
//   - Initial primary server establishment
//   - Backup server promotion
//   - Primary failure handling
//   - Server restart scenarios
//   - View acknowledgment requirements
//   - Uninitialized server restrictions
func Test1(t *testing.T) {
	// Set maximum number of OS threads to ensure proper concurrency testing
	runtime.GOMAXPROCS(4)

	// Create view service and client clerks
	vshost := port("v")
	vs := StartServer(vshost)
	defer vs.Kill() // Ensure cleanup

	ck1 := MakeClerk(port("1"), vshost)
	ck2 := MakeClerk(port("2"), vshost)
	ck3 := MakeClerk(port("3"), vshost)

	// Verify no primary exists initially
	if ck1.Primary() != "" {
		t.Fatalf("there was a primary too soon")
	}

	// Test 1: Establish the very first primary server
	fmt.Printf("Test: First primary ...\n")
	testFirstPrimary(t, ck1)
	fmt.Printf("  ... Passed\n")

	// Test 2: Establish the first backup server
	fmt.Printf("Test: First backup ...\n")
	testFirstBackup(t, ck1, ck2)
	fmt.Printf("  ... Passed\n")

	// Test 3: Primary failure and backup takeover
	fmt.Printf("Test: Backup takes over if primary fails ...\n")
	testBackupTakeover(t, ck1, ck2)
	fmt.Printf("  ... Passed\n")

	// Test 4: Restarted server becomes backup
	fmt.Printf("Test: Restarted server becomes backup ...\n")
	testRestartedServerBecomesBackup(t, ck1, ck2)
	fmt.Printf("  ... Passed\n")

	// Test 5: Idle third server becomes backup when primary fails
	fmt.Printf("Test: Idle third server becomes backup if primary fails ...\n")
	testIdleServerBecomesBackup(t, ck1, ck2, ck3)
	fmt.Printf("  ... Passed\n")

	// Test 6: Restarted primary treated as dead
	fmt.Printf("Test: Restarted primary treated as dead ...\n")
	testRestartedPrimaryTreatedAsDead(t, ck1, ck3)
	fmt.Printf("  ... Passed\n")

	// Test 7: View service waits for primary acknowledgment
	fmt.Printf("Test: Viewserver waits for primary to ack view ...\n")
	testViewserverWaitsForAck(t, ck1, ck2, ck3)
	fmt.Printf("  ... Passed\n")

	// Test 8: Uninitialized server cannot become primary
	fmt.Printf("Test: Uninitialized server can't become primary ...\n")
	testUninitializedServerCannotBecomePrimary(t, ck1, ck2, ck3)
	fmt.Printf("  ... Passed\n")
}

// testFirstPrimary tests the establishment of the very first primary server.
// It verifies that when a server pings with viewnum=0, it becomes the primary
// in view 1 if no primary currently exists.
func testFirstPrimary(t *testing.T, ck1 *Clerk) {
	for i := 0; i < DeadPings*2; i++ {
		view, err := ck1.Ping(0)
		if err != nil {
			t.Fatalf("Ping failed: %v", err)
		}
		if view.Primary == ck1.me {
			break
		}
		time.Sleep(PingInterval)
	}
	check(t, ck1, ck1.me, "", 1)
}

// testFirstBackup tests the establishment of the first backup server.
// It verifies that when a second server pings with viewnum=0 after the primary
// has acknowledged its view, it becomes the backup in a new view.
func testFirstBackup(t *testing.T, ck1, ck2 *Clerk) {
	vx, ok := ck1.Get()
	if !ok {
		t.Fatalf("failed to get initial view")
	}
	
	for i := 0; i < DeadPings*2; i++ {
		// Primary acknowledges current view
		_, err := ck1.Ping(1)
		if err != nil {
			t.Fatalf("Primary ping failed: %v", err)
		}
		
		// Second server volunteers to be backup
		view, err := ck2.Ping(0)
		if err != nil {
			t.Fatalf("Backup ping failed: %v", err)
		}
		if view.Backup == ck2.me {
			break
		}
		time.Sleep(PingInterval)
	}
	check(t, ck1, ck1.me, ck2.me, vx.Viewnum+1)
}

// testBackupTakeover tests that when the primary fails, the backup takes over.
// It simulates primary failure by stopping primary pings and verifies that
// the backup becomes the new primary with no backup.
func testBackupTakeover(t *testing.T, ck1, ck2 *Clerk) {
	// Primary acknowledges current view
	_, err := ck1.Ping(2)
	if err != nil {
		t.Fatalf("Primary ping failed: %v", err)
	}
	
	// Backup acknowledges current view
	vx, err := ck2.Ping(2)
	if err != nil {
		t.Fatalf("Backup ping failed: %v", err)
	}
	
	// Wait for primary to be declared dead and backup to take over
	for i := 0; i < DeadPings*2; i++ {
		v, err := ck2.Ping(vx.Viewnum)
		if err != nil {
			t.Fatalf("Backup ping failed: %v", err)
		}
		if v.Primary == ck2.me && v.Backup == "" {
			break
		}
		time.Sleep(PingInterval)
	}
	check(t, ck2, ck2.me, "", vx.Viewnum+1)
}

// testRestartedServerBecomesBackup tests that a restarted server can become backup.
// It simulates a server restart (ping with viewnum=0) and verifies it becomes
// the backup when the current view is acknowledged.
func testRestartedServerBecomesBackup(t *testing.T, ck1, ck2 *Clerk) {
	vx, ok := ck2.Get()
	if !ok {
		t.Fatalf("failed to get current view")
	}
	
	// Current primary acknowledges view
	_, err := ck2.Ping(vx.Viewnum)
	if err != nil {
		t.Fatalf("Primary ping failed: %v", err)
	}
	
	// Restarted server volunteers to be backup
	for i := 0; i < DeadPings*2; i++ {
		_, err := ck1.Ping(0)
		if err != nil {
			t.Fatalf("Restarted server ping failed: %v", err)
		}
		
		v, err := ck2.Ping(vx.Viewnum)
		if err != nil {
			t.Fatalf("Primary ping failed: %v", err)
		}
		if v.Primary == ck2.me && v.Backup == ck1.me {
			break
		}
		time.Sleep(PingInterval)
	}
	check(t, ck2, ck2.me, ck1.me, vx.Viewnum+1)
}

// testIdleServerBecomesBackup tests that an idle third server becomes backup
// when the primary fails. It verifies the promotion chain: primary fails,
// backup becomes primary, idle server becomes backup.
func testIdleServerBecomesBackup(t *testing.T, ck1, ck2, ck3 *Clerk) {
	vx, ok := ck2.Get()
	if !ok {
		t.Fatalf("failed to get current view")
	}
	
	// Current primary acknowledges view
	_, err := ck2.Ping(vx.Viewnum)
	if err != nil {
		t.Fatalf("Primary ping failed: %v", err)
	}
	
	// Wait for primary failure and promotion chain
	for i := 0; i < DeadPings*2; i++ {
		// Third server volunteers to be backup
		_, err := ck3.Ping(0)
		if err != nil {
			t.Fatalf("Third server ping failed: %v", err)
		}
		
		// Backup should become primary, third server should become backup
		v, err := ck1.Ping(vx.Viewnum)
		if err != nil {
			t.Fatalf("Backup ping failed: %v", err)
		}
		if v.Primary == ck1.me && v.Backup == ck3.me {
			break
		}
		vx = v
		time.Sleep(PingInterval)
	}
	check(t, ck1, ck1.me, ck3.me, vx.Viewnum+1)
}

// testRestartedPrimaryTreatedAsDead tests that a primary that restarts
// (pings with viewnum=0) is treated as dead and replaced.
func testRestartedPrimaryTreatedAsDead(t *testing.T, ck1, ck3 *Clerk) {
	vx, ok := ck1.Get()
	if !ok {
		t.Fatalf("failed to get current view")
	}
	
	// Primary acknowledges current view
	_, err := ck1.Ping(vx.Viewnum)
	if err != nil {
		t.Fatalf("Primary ping failed: %v", err)
	}
	
	// Primary restarts (pings with viewnum=0) - should be treated as dead
	for i := 0; i < DeadPings*2; i++ {
		_, err := ck1.Ping(0)
		if err != nil {
			t.Fatalf("Restarted primary ping failed: %v", err)
		}
		
		_, err = ck3.Ping(vx.Viewnum)
		if err != nil {
			t.Fatalf("Backup ping failed: %v", err)
		}
		
		v, ok := ck3.Get()
		if !ok {
			t.Fatalf("failed to get view")
		}
		if v.Primary != ck1.me {
			break
		}
		time.Sleep(PingInterval)
	}
	
	vy, ok := ck3.Get()
	if !ok {
		t.Fatalf("failed to get final view")
	}
	if vy.Primary != ck3.me {
		t.Fatalf("expected primary=%v, got %v", ck3.me, vy.Primary)
	}
}

// testViewserverWaitsForAck tests that the view service waits for primary
// acknowledgment before proceeding to new views. It verifies that if a primary
// fails before acknowledging a view, the view service cannot make progress.
func testViewserverWaitsForAck(t *testing.T, ck1, ck2, ck3 *Clerk) {
	// Set up a view with ck3 as primary and ck1 as backup, but don't ack
	vx, ok := ck1.Get()
	if !ok {
		t.Fatalf("failed to get initial view")
	}
	
	// Establish new view with ck3 as primary, ck1 as backup
	for i := 0; i < DeadPings*3; i++ {
		_, err := ck1.Ping(0)
		if err != nil {
			t.Fatalf("Server ping failed: %v", err)
		}
		
		_, err = ck3.Ping(vx.Viewnum)
		if err != nil {
			t.Fatalf("Primary ping failed: %v", err)
		}
		
		v, ok := ck1.Get()
		if !ok {
			t.Fatalf("failed to get view")
		}
		if v.Viewnum > vx.Viewnum {
			break
		}
		time.Sleep(PingInterval)
	}
	check(t, ck1, ck3.me, ck1.me, vx.Viewnum+1)
	
	// ck3 is primary but never acknowledged the view
	// Let ck3 die and verify ck1 is not promoted (view service waits for ack)
	vy, ok := ck1.Get()
	if !ok {
		t.Fatalf("failed to get view")
	}
	
	for i := 0; i < DeadPings*3; i++ {
		v, err := ck1.Ping(vy.Viewnum)
		if err != nil {
			t.Fatalf("Backup ping failed: %v", err)
		}
		if v.Viewnum > vy.Viewnum {
			break
		}
		time.Sleep(PingInterval)
	}
	check(t, ck2, ck3.me, ck1.me, vy.Viewnum)
}

// testUninitializedServerCannotBecomePrimary tests that an uninitialized server
// (one that has never been primary or backup) cannot become primary even if
// all other servers are dead.
func testUninitializedServerCannotBecomePrimary(t *testing.T, ck1, ck2, ck3 *Clerk) {
	// Ensure all servers are in a known state
	for i := 0; i < DeadPings*2; i++ {
		v, ok := ck1.Get()
		if !ok {
			t.Fatalf("failed to get view")
		}
		
		_, err := ck1.Ping(v.Viewnum)
		if err != nil {
			t.Fatalf("Server ping failed: %v", err)
		}
		
		_, err = ck2.Ping(0)
		if err != nil {
			t.Fatalf("Server ping failed: %v", err)
		}
		
		_, err = ck3.Ping(v.Viewnum)
		if err != nil {
			t.Fatalf("Server ping failed: %v", err)
		}
		time.Sleep(PingInterval)
	}
	
	// Let ck2 (uninitialized server) ping for a while
	for i := 0; i < DeadPings*2; i++ {
		_, err := ck2.Ping(0)
		if err != nil {
			t.Fatalf("Uninitialized server ping failed: %v", err)
		}
		time.Sleep(PingInterval)
	}
	
	// Verify ck2 did not become primary
	vz, ok := ck2.Get()
	if !ok {
		t.Fatalf("failed to get final view")
	}
	if vz.Primary == ck2.me {
		t.Fatalf("uninitialized backup promoted to primary")
	}
}
