package numbers

import "testing"

func TestDecimalToSatoshis(t *testing.T) {
	assertSatEquals := func(expected string, input string) {
		actual, err := DecimalToSatoshis(input)
		if err != nil {
			t.Error(err)
		}
		if expected != actual {
			t.Errorf("expected %s, got %s, input %s", expected, actual, input)
		}
	}

	assertSatError := func(input string) {
		actual, err := DecimalToSatoshis(input)
		if err == nil {
			t.Errorf("Expected error but no error: got %s, input %s", actual, input)
		}
	}

	assertSatEquals("10", "1.0")
	assertSatEquals("1", "0.1")
	assertSatEquals("13602", "136.02")
	assertSatEquals("13602", "0136.02")
	assertSatEquals("1500000", "0.01500000")
	assertSatEquals("0", "0")
	assertSatEquals("2030", "0.002030")
	assertSatEquals("101010", "0101010")
	assertSatEquals("11001100", "0011001100")
	assertSatEquals("376", " 376")
	assertSatEquals("376", "376 ")

	assertSatError("12NotNumber34")
	assertSatError("12,34")
	assertSatError("")
	assertSatError(" ")
	assertSatError("37 6")
	assertSatError("37,6")
}
