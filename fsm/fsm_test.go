package fsm

import (
	"testing"

	"github.com/go-redis/redis/v8"
)

func testEq(a, b []string) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFSM(t *testing.T) {
	path := "../examples/00_test/"
	domain := Create(&path)
	commandList := []string{"turn_on", "turn_off", "hello_universe"}
	if !testEq(domain.CommandList, commandList) {
		t.Errorf("domain.CommandList is incorrect, got: %v, want: %v.", domain.CommandList, commandList)
	}

	machine := FSM{State: 0}
	extension := LoadExtensions(&path)
	// if err != nil {
	// 	t.Errorf(err.Error())
	// }

	resp1 := machine.ExecuteCmd("turn_on", "turn_on", domain, extension)
	if resp1 != "Turning on." {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp1, "Turning on.")
	}

	resp2 := machine.ExecuteCmd("turn_on", "turn_on", domain, extension)
	if resp2 != "Can't do that." {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp2, "Can't do that.")
	}

	resp3 := machine.ExecuteCmd("hello_universe", "hello", domain, extension)
	if resp3 != "Hello Universe" {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp3, "Hello Universe")
	}

	// testFuncExists := extension.GetFunc("ext_any")
	// testFuncExistsNot := extension.GetFunc("ext_any_other")
	// if testFuncExists == nil {
	// 	t.Errorf("GetFunc is incorrect, 'ext_any' should exist.")
	// }
	// if testFuncExistsNot != nil {
	// 	t.Errorf("GetFunc is incorrect, 'ext_any_other' should not exist.")
	// }
}

func TestCacheStore(t *testing.T) {
	machines := &CacheStoreFSM{}
	if resp1 := machines.Exists("foo"); resp1 != false {
		t.Errorf("incorrect, got: %v, want: %v.", resp1, "false")
	}

	machines.Set(
		"foo",
		&FSM{
			State: 0,
			Slots: make(map[string]string),
		},
	)
	if resp2 := machines.Get("foo"); resp2.State != 0 {
		t.Errorf("incorrect, got: %v, want: %v.", resp2, "0")
	}

	newFsm := &FSM{
		State: 1,
		Slots: map[string]string{
			"abc": "xyz",
		},
	}
	machines.Set("foo", newFsm)
	if resp3 := machines.Get("foo"); resp3.State != 1 {
		t.Errorf("incorrect, got: %v, want: %v.", resp3, "1")
	}
}

func TestRedisStore(t *testing.T) {
	var rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "pass",
		DB:       0,
	})

	machines := &RedisStoreFSM{R: rdb}
	if resp1 := machines.Exists("foo"); resp1 != false {
		t.Errorf("incorrect, got: %v, want: %v.", resp1, "false")
	}

	machines.Set(
		"foo",
		&FSM{
			State: 0,
			Slots: make(map[string]string),
		},
	)
	if resp2 := machines.Get("foo"); resp2.State != 0 {
		t.Errorf("incorrect, got: %v, want: %v.", resp2, "0")
	}

	newFsm := &FSM{
		State: 1,
		Slots: map[string]string{
			"abc": "xyz",
		},
	}
	machines.Set("foo", newFsm)
	if resp3 := machines.Get("foo"); resp3.State != 1 {
		t.Errorf("incorrect, got: %v, want: %v.", resp3, "1")
	}
}

func TestExt(t *testing.T) {
	greetFunc := func(req *Request) (res *Response) {
		return &Response{
			FSM: req.FSM,
			Res: "Hello Universe",
		}
	}

	var myExtMap = ExtensionMap{
		"ext_any": greetFunc,
	}

	listener := &ListenerRPC{ExtensionMap: myExtMap}

	req1 := &Request{
		FSM: &FSM{
			State: 0,
			Slots: make(map[string]string),
		},
		Req: "ext_any",
	}
	res1 := new(Response)
	listener.GetFunc(req1, res1)
	if res1.Res != "Hello Universe" {
		t.Errorf("incorrect, got: %v, want: %v.", res1.Res, "Hello Universe")
	}

	req2 := new(Request)
	res2 := new(GetAllFuncsResponse)
	listener.GetAllFuncs(req2, res2)
	if len(res2.Res) != 1 {
		t.Errorf("incorrect, got: %v, want: %v.", len(res2.Res), "1")
	}
}
