package tools

import (
	"sync"
	"testing"
)

func TestNewSnowflake(t *testing.T) {
	var i, j int64
	for i = 0; i < 32; i++ {
		for j = 0; j < 32; j++ {
			_, err := NewSnowflake(i, j)
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestNextVal(t *testing.T) {
	s, err := NewSnowflake(0, 0)
	if err != nil {
		t.Error(err)
	}
	var i int64
	for i = 0; i < sequenceMask*10; i++ {
		val, err := s.NextVal()
		if val == 0 || err != nil {
			t.Fail()
		}
	}
}

func TestUnique(t *testing.T) {
	var wg sync.WaitGroup
	var check sync.Map
	s, err := NewSnowflake(0, 0)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Add(-1)
			val, err := s.NextVal()
			if err != nil {
				t.Fail()
				return
			}
			if _, ok := check.Load(val); ok {
				t.Fail()
				return
			}
			check.Store(val, 0)
			if val == 0 {
				t.Fail()
				return
			}
		}()
	}
	wg.Wait()
}

func TestGetTime(t *testing.T) {
	s, err := NewSnowflake(0, 1)
	if err != nil {
		t.Error(err)
	}
	val, err := s.NextVal()
	if err != nil {
		t.Error(err)
	}
	t.Logf("time:%v", GetGenTime(val))
}

func TestGetDeviceID(t *testing.T) {
	s, err := NewSnowflake(28, 11)
	if err != nil {
		t.Error(err)
	}
	val, err := s.NextVal()
	if err != nil {
		t.Error(err)
	}
	datacenterid, workerid := GetDeviceID(val)
	if datacenterid != 28 || workerid != 11 {
		t.Fail()
	}
}
