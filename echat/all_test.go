package main

import "testing"

func TestTest (t *testing.T) {
  if 1 != 1  {
    t.Errorf("Error: 1 does not equal 1")
  }
}
