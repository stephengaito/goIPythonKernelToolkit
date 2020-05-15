// +build cGoTests

package goIPyRubyAdaptor

import(
  "os/exec"
  "testing"
)

func Test_IPyRubyData(t *testing.T) {
  t.Logf("IPyRubyData tests started\n")
  cmd := exec.Command("./runIPyRubyDataTests")
  cmdOut, err := cmd.Output()
  t.Logf("%s", cmdOut)
  if err != nil {
    t.Errorf("IPyRubyData test failed ERROR: %s", err)
  }
  t.Logf("IPyRubyData tests done\n")
}
