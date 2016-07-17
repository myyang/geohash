Geohash
=======

.. image:: https://travis-ci.org/myyang/geohash.svg?branch=master
    :target: https://travis-ci.org/myyang/geohash


* Geohash
* Geohash-36

Example
-------

Install package:

.. code:: shell 

    go get github.com/myyang/geohash


Example:

.. code:: golang

    import (
        "github.com/myyang/geohash"
        "fmt"
    )

    func main() {
        cryptor := NewDefaultGeoHash()
        hashValue := cryptor.Encode(12.04512315, 118.20385763, 9)
        fmt.Println(hashValue)  // print "wdhh9b9rv"
    }

TODO
----

* [ ] more invalid checking & testing
