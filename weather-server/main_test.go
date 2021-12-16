package main

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHeadIndex(t *testing.T) {
	//0.5 * ($t + 61.0 + (($t - 68.0) * 1.2) + ($rh * 0.094))
    //-42.379 + 2.04901523 * $t + 10.1433127*$rh - .22475441*$t*$rh - .00683783 *$t * $t - .05481717 * $rh * $rh + .00122874*$t*$t*$rh + .00085282 *$t * $rh *$rh - .00000199 *$t *$t *$rh * $rh;$a = 0;if ($rh < 13 && ($t >= 80 && $t <= 112)) {$a=((13 - $rh ) / 4) * sqrt((17-abs($t - 95))/17);$a = -$a;};if ($rh > 85 && ($t >= 80 && $t <= 87)) {$a=(($rh - 85)/10) * ((87 - $t) / 5)


}