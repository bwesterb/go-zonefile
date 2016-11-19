package zonefile_test

import (
	"bytes"
	"fmt"
	"github.com/bwesterb/go-zonefile"
	"testing"
)

// Loading and saving a zonefile shouldn't do anything
func TestLoadThenSave(t *testing.T) {
	for i, test := range tests {
		z, e := zonefile.Load([]byte(test))
		if e != nil {
			t.Fatal(i, "error loading:", e.LineNo(), e)
		}
		if !bytes.Equal(z.Save(), []byte(test)) {
			t.Fatal("Save o Load != identity")
		}
	}
}

func ExampleLoad() {
	zf, err := zonefile.Load([]byte(
		"@	IN	SOA	NS1.NAMESERVER.NET.	HOSTMASTER.MYDOMAIN.COM.	(\n" +
			"            1406291485	 ;serial\n" +
			"            3600	 ;refresh\n" +
			"            600	 ;retry\n" +
			"            604800	 ;expire\n" +
			"            86400	 ;minimum ttl\n" +
			")\n" +
			"\n" +
			"@	NS	NS1.NAMESERVER.NET.\n" +
			"@	NS	NS2.NAMESERVER.NET.\n"))
	if err != nil {
		fmt.Println("Parsing error", err, "on line", err.LineNo())
		return
	}
	fmt.Println(zf)
	// Output: <Zonefile with 3 entries>
}

var tests = [...]string{`$ORIGIN MYDOMAIN.COM.
$TTL 3600
@	IN	SOA	NS1.NAMESERVER.NET.	HOSTMASTER.MYDOMAIN.COM.	(
			1406291485	 ;serial
			3600	 ;refresh
			600	 ;retry
			604800	 ;expire
			86400	 ;minimum ttl
)

@	NS	NS1.NAMESERVER.NET.
@	NS	NS2.NAMESERVER.NET.

@	MX	0	mail1
@	MX	10	mail2

	A	1.1.1.1
@	A	127.0.0.1
www	A	127.0.0.1
mail	A	127.0.0.1
			A 1.2.3.4
tst 300 IN A 101.228.10.127;this is a comment

@	AAAA	::1
mail	AAAA	2001:db8::1

mail1	CNAME	mail
mail2	CNAME	mail

treefrog.ca. IN TXT "v=spf1 a mx a:mail.treefrog.ca a:webmail.treefrog.ca ip4:76.75.250.33 ?all"
treemonkey.ca. IN TXT "v=DKIM1\; k=rsa\; p=MIGf..."`,
	`$ORIGIN 0.168.192.IN-ADDR.ARPA.
$TTL 3600
@	IN	SOA	NS1.NAMESERVER.NET.	HOSTMASTER.MYDOMAIN.COM.	(
			1406291485	 ;serial
			3600	 ;refresh
			600	 ;retry
			604800	 ;expire
			86400	 ;minimum ttl
)

@	NS	NS1.NAMESERVER.NET.
@	NS	NS2.NAMESERVER.NET.

1	PTR	HOST1.MYDOMAIN.COM.
2	PTR	HOST2.MYDOMAIN.COM.

$ORIGIN 30.168.192.in-addr.arpa.
3	PTR	HOST3.MYDOMAIN.COM.
4	PTR	HOST4.MYDOMAIN.COM.
	PTR HOST5.MYDOMAIN.COM.

$ORIGIN 168.192.in-addr.arpa.
10.3	PTR	HOST3.MYDOMAIN.COM.
10.4	PTR	HOST4.MYDOMAIN.COM.`,
	`$ORIGIN 0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa.
$TTL 3600
@	IN	SOA	NS1.NAMESERVER.NET.	HOSTMASTER.MYDOMAIN.COM.	(
			1406291485	 ;serial
			3600	 ;refresh
			600	 ;retry
			604800	 ;expire
			86400	 ;minimum ttl
)

@	NS	NS1.NAMESERVER.NET.
@	NS	NS2.NAMESERVER.NET.

1	PTR	HOST1.MYDOMAIN.COM.
2	PTR	HOST2.MYDOMAIN.COM.`,
	`$ORIGIN example.com.     ; designates the start of this zone file in the namespace
$TTL 1h                  ; default expiration time of all resource records without their own TTL value
example.com.  IN  SOA   ns.example.com. username.example.com. ( 2007120710 1d 2h 4w 1h )
example.com.  IN  NS    ns                    ; ns.example.com is a nameserver for example.com
example.com.  IN  NS    ns.somewhere.example. ; ns.somewhere.example is a backup nameserver for example.com
example.com.  IN  MX    10 mail.example.com.  ; mail.example.com is the mailserver for example.com
@             IN  MX    20 mail2.example.com. ; equivalent to above line, "@" represents zone origin
@             IN  MX    50 mail3              ; equivalent to above line, but using a relative host name
example.com.  IN  A     192.0.2.1             ; IPv4 address for example.com
              IN  AAAA  2001:db8:10::1        ; IPv6 address for example.com
ns            IN  A     192.0.2.2             ; IPv4 address for ns.example.com
              IN  AAAA  2001:db8:10::2        ; IPv6 address for ns.example.com
www           IN  CNAME example.com.          ; www.example.com is an alias for example.com
wwwtest       IN  CNAME www                   ; wwwtest.example.com is another alias for www.example.com
mail          IN  A     192.0.2.3             ; IPv4 address for mail.example.com
mail2         IN  A     192.0.2.4             ; IPv4 address for mail2.example.com
mail3         IN  A     192.0.2.5             ; IPv4 address for mail3.example.com`}
