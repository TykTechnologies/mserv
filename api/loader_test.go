package api

import "testing"

func TestPluginLoadAndDispatch(t *testing.T) {
	//d := Dispatcher{
	//	store: RuntimeStore{},
	//}

	_, err := LoadPlugin("MyPreHook", "plugins/", "plugin.so")
	if err != nil {
		t.Fatal(err)
	}

}
