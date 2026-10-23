package cache

import "testing"

func TestPutGetEvict(t *testing.T) {
	c := New(2)
	c.Put("a", "1")
	c.Put("b", "2")
	c.Put("c", "3")
	if _, ok := c.Get("a"); ok {
		t.Fatal("a should be evicted")
	}
	if v, ok := c.Get("c"); !ok || v != "3" || c.Len() != 2 {
		t.Fatalf("c = %q %v len %d", v, ok, c.Len())
	}
}

func TestGetOrLoad(t *testing.T) {
	c := New(2)
	v, err := c.GetOrLoad("k", func() (string, error) { return "v", nil })
	if err != nil || v != "v" {
		t.Fatalf("GetOrLoad = %q, %v", v, err)
	}
}
