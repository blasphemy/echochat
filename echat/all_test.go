package main

import "testing"

var 

func TestTest (t *testing.T) {
  if 1 != 1  {
    t.Errorf("Error: 1 does not equal 1")
  }
}
func TestAddUserToList(t *testing.T) {
  userlist = make(map[int]User)
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
