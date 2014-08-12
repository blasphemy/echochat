package main

import "testing"

func SetupTest() {
  userlist = make(map[int]User)
}
func TestSetupTest(t *testing.T) {
SetupTest()
  if userlist == nil {
   t.Errorf("userlist is nil")
  }
}

func TestTest (t *testing.T) {
  if 1 != 1  {
    t.Errorf("Error: 1 does not equal 1")
  }
}
func TestAddUserToList(t *testing.T) {
  user := User{id: 1}
  AddUserToList(user)
  if len(userlist) != 1 {
    t.Errorf("Error, adding user to list did not increse length")
  }
}

func TestRemoveUser(t *testing.T) {
  if len(userlist) != 1 {
    t.Errorf("Error, list has not been initialized")
  }
}
