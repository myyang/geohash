/* Package geohash provides cryptors which implement geohash algorithm
Currently implemented:

* Geohash
* Geohash-36

Let's see an example:

   import (
       "github.com/myyang/geohash"
       "fmt"
   )

   func main() {
       cryptor := NewDefaultGeoHash()
       hashValue := cryptor.Encode(12.04512315, 118.20385763, 9)
       fmt.Println(hashValue)  // print "wdhh9b9rv"
   }

*/
package geohash
