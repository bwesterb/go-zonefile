package zonefile_test

import (
	"bytes"
	"fmt"
	"github.com/bwesterb/go-zonefile"
	"testing"
)

// XXX more tests for AddEntry.
// XXX test Value() and SetValue()

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

func TestSetAttributes(t *testing.T) {
	zf, err := zonefile.Load([]byte(" IN A 1.2.3.4"))
	if err != nil {
		t.Fatal("Couldn't parse simple zonefile")
	}
	zf.Entries()[0].SetDomain([]byte("test"))
	if !bytes.Equal(zf.Save(), []byte("test IN A 1.2.3.4")) {
		t.Fatal("Setting domain failed")
	}
	zf.Entries()[0].SetDomain([]byte("test2"))
	if !bytes.Equal(zf.Save(), []byte("test2 IN A 1.2.3.4")) {
		t.Fatal("Setting domain failed")
	}
	zf.Entries()[0].SetDomain([]byte(""))
	if !bytes.Equal(zf.Save(), []byte(" IN A 1.2.3.4")) {
		t.Fatal("Setting domain failed")
	}
	zf.Entries()[0].SetClass([]byte(""))
	if !bytes.Equal(zf.Save(), []byte("  A 1.2.3.4")) {
		t.Fatal("Setting class failed")
	}
	zf.Entries()[0].SetClass([]byte("IN"))
	if !bytes.Equal(zf.Save(), []byte(" IN A 1.2.3.4")) {
		t.Fatal("Setting class failed")
	}
	zf.Entries()[0].SetDomain([]byte("test4"))
	if !bytes.Equal(zf.Save(), []byte("test4 IN A 1.2.3.4")) {
		t.Fatal("Setting class failed")
	}
	zf.Entries()[0].SetClass([]byte(""))
	if !bytes.Equal(zf.Save(), []byte("test4  A 1.2.3.4")) {
		t.Fatal("Setting class failed")
	}
	zf.Entries()[0].SetTTL(12)
	if !bytes.Equal(zf.Save(), []byte("test4 12  A 1.2.3.4")) {
		t.Fatal("Setting class failed")
	}
	zf.Entries()[0].SetTTL(14)
	if !bytes.Equal(zf.Save(), []byte("test4 14  A 1.2.3.4")) {
		t.Fatal("Setting class failed")
	}
	if *zf.Entries()[0].TTL() != 14 {
		t.Fatal("TTL wasn't properly set")
	}
	zf.Entries()[0].RemoveTTL()
	if !bytes.Equal(zf.Save(), []byte("test4   A 1.2.3.4")) {
		t.Fatal("Setting class failed")
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
	fmt.Println(len(zf.Entries()))
	// Output: 3
}

func ExampleParseEntry() {
	entry, err := zonefile.ParseEntry([]byte(" IN MX 100 alpha.example.com."))
	if err != nil {
		fmt.Println("Parsing error", err, "on line", err.LineNo())
		return
	}
	fmt.Println(entry)
	// Output: <Entry dom="" ttl="" cls="IN" typ="MX" ["100" "alpha.example.com."]>
}

func ExampleNew() {
	z := zonefile.New()
	z.AddA("", "3.2.3.2")
	z.AddA("www", "1.2.3.4")
	z.AddA("irc", "2.2.2.2").SetTTL(12)
	fmt.Println(z)
	// Output: <Zonefile [<Entry dom="" ttl="" cls="" typ="" ["3.2.3.2"]> <Entry dom="www" ttl="" cls="" typ="" ["1.2.3.4"]> <Entry dom="irc" ttl="12" cls="" typ="" ["2.2.2.2"]>]>
}

func ExampleZonefile_AddA() {
	z := zonefile.New()
	z.AddA("", "3.2.3.2")
	z.AddA("www", "1.2.3.4")
	z.AddA("irc", "2.2.2.2").SetTTL(12)
	fmt.Println(z)
	// Output: <Zonefile [<Entry dom="" ttl="" cls="" typ="" ["3.2.3.2"]> <Entry dom="www" ttl="" cls="" typ="" ["1.2.3.4"]> <Entry dom="irc" ttl="12" cls="" typ="" ["2.2.2.2"]>]>
}

func ExampleZonefile_AddEntry() {
	z := zonefile.New()
	entry, _ := zonefile.ParseEntry([]byte("irc IN A 1.2.3.4"))
	z.AddEntry(entry)
	fmt.Println(z)
	// Output: <Zonefile [<Entry dom="irc" ttl="" cls="IN" typ="A" ["1.2.3.4"]>]>
}

func ExampleZonefile_Save() {
	z := zonefile.New()
	entry, _ := zonefile.ParseEntry([]byte("irc IN A 1.2.3.4"))
	z.AddEntry(entry)
	entry, _ = zonefile.ParseEntry([]byte("www IN A 2.1.4.3"))
	z.AddEntry(entry)
	fmt.Println(string(z.Save()))
	// Output: irc IN A 1.2.3.4
	// www IN A 2.1.4.3
}

func ExampleZonefile_Entries() {
	zf, err := zonefile.Load([]byte(`
$TTL 3600
@	IN	SOA	NS1.NAMESERVER.NET.	HOSTMASTER.MYDOMAIN.COM.	(
			1406291485	 ;serial
			3600	 ;refresh
			600	 ;retry
			604800	 ;expire
			86400	 ;minimum ttl
)

	A	1.1.1.1
@	A	127.0.0.1
www	A	127.0.0.1
mail	A	127.0.0.1
			A 1.2.3.4
tst 300 IN A 101.228.10.127;this is a comment`))
	if err != nil {
		fmt.Println("Error parsing zonefile:", err)
		return
	}
	for _, e := range zf.Entries() {
		fmt.Println(e)
	}
	// Output: <Entry cmd="$TTL" ["3600"]>
	// <Entry dom="@" ttl="" cls="IN" typ="SOA" ["NS1.NAMESERVER.NET." "HOSTMASTER.MYDOMAIN.COM." "1406291485" "3600" "600" "604800" "86400"]>
	// <Entry dom="" ttl="" cls="" typ="A" ["1.1.1.1"]>
	// <Entry dom="@" ttl="" cls="" typ="A" ["127.0.0.1"]>
	// <Entry dom="www" ttl="" cls="" typ="A" ["127.0.0.1"]>
	// <Entry dom="mail" ttl="" cls="" typ="A" ["127.0.0.1"]>
	// <Entry dom="" ttl="" cls="" typ="A" ["1.2.3.4"]>
	// <Entry dom="tst" ttl="300" cls="IN" typ="A" ["101.228.10.127"]>
}

func ExampleEntry_SetValue() {
	entry, _ := zonefile.ParseEntry([]byte("irc IN A 1.2.3.4"))
	fmt.Println(entry)
	entry.SetValue(0, []byte("4.3.2.1"))
	fmt.Println(entry)
	// Output: <Entry dom="irc" ttl="" cls="IN" typ="A" ["1.2.3.4"]>
	// <Entry dom="irc" ttl="" cls="IN" typ="A" ["4.3.2.1"]>
}

func ExampleEntry_RemoveTTL() {
	entry, _ := zonefile.ParseEntry([]byte("irc 12 IN A 1.2.3.4"))
	fmt.Println(entry)
	entry.RemoveTTL()
	fmt.Println(entry)
	// Output: <Entry dom="irc" ttl="12" cls="IN" typ="A" ["1.2.3.4"]>
	// <Entry dom="irc" ttl="" cls="IN" typ="A" ["1.2.3.4"]>
}

func ExampleEntry_SetTTL() {
	entry, _ := zonefile.ParseEntry([]byte("irc 12 IN A 1.2.3.4"))
	fmt.Println(entry)
	entry.SetTTL(14)
	fmt.Println(entry)
	// Output: <Entry dom="irc" ttl="12" cls="IN" typ="A" ["1.2.3.4"]>
	// <Entry dom="irc" ttl="14" cls="IN" typ="A" ["1.2.3.4"]>
}

func ExampleEntry_Domain() {
	entry, _ := zonefile.ParseEntry([]byte("irc IN A 1.2.3.4"))
	fmt.Printf("%q\n", entry.Domain())
	entry, _ = zonefile.ParseEntry([]byte(" IN A 4.3.2.1"))
	fmt.Printf("%q\n", entry.Domain())
	// Output: "irc"
	// ""
}

func ExampleEntry_Class() {
	entry, _ := zonefile.ParseEntry([]byte("irc A 1.2.3.4"))
	fmt.Printf("%q\n", entry.Class())
	entry, _ = zonefile.ParseEntry([]byte("irc IN A 4.3.2.1"))
	fmt.Printf("%q\n", entry.Class())
	// Output: ""
	// "IN"
}

func ExampleEntry_Type() {
	entry, _ := zonefile.ParseEntry([]byte("irc A 1.2.3.4"))
	fmt.Printf("%q\n", entry.Type())
	entry, _ = zonefile.ParseEntry([]byte("irc AAAA ::1"))
	fmt.Printf("%q\n", entry.Type())
	// Output: "A"
	// "AAAA"
}

func ExampleEntry_TTL() {
	entry, _ := zonefile.ParseEntry([]byte("irc A 1.2.3.4"))
	fmt.Printf("%v\n", entry.TTL() == nil)
	entry, _ = zonefile.ParseEntry([]byte("irc 12 A 1.2.3.4"))
	fmt.Printf("%v\n", *entry.TTL())
	// Output: true
	// 12
}

func ExampleEntry_Command() {
	entry, _ := zonefile.ParseEntry([]byte("$TTL 123"))
	fmt.Printf("%q\n", entry.Command())
	// Output: "$TTL"
}

func ExampleEntry_Values() {
	entry, _ := zonefile.ParseEntry([]byte(`
@	IN	SOA	NS1.NAMESERVER.NET.	HOSTMASTER.MYDOMAIN.COM.	(
			1406291485	 ;serial
			3600	 ;refresh
			600	 ;retry
			604800	 ;expire
			86400	 ;minimum ttl
)`))
	fmt.Printf("%q\n", entry.Values())
	// Output: ["NS1.NAMESERVER.NET." "HOSTMASTER.MYDOMAIN.COM." "1406291485" "3600" "600" "604800" "86400"]
}

func ExampleEntry_SetDomain() {
	entry, _ := zonefile.ParseEntry([]byte("irc IN A 1.2.3.4"))
	fmt.Println(entry)
	entry.SetDomain([]byte(""))
	fmt.Println(entry)
	entry.SetDomain([]byte("chat"))
	fmt.Println(entry)
	// Output: <Entry dom="irc" ttl="" cls="IN" typ="A" ["1.2.3.4"]>
	// <Entry dom="" ttl="" cls="IN" typ="A" ["1.2.3.4"]>
	// <Entry dom="chat" ttl="" cls="IN" typ="A" ["1.2.3.4"]>
}

func ExampleEntry_SetClass() {
	entry, _ := zonefile.ParseEntry([]byte("irc IN A 1.2.3.4"))
	fmt.Println(entry)
	entry.SetClass([]byte(""))
	fmt.Println(entry)
	entry.SetClass([]byte("IN"))
	fmt.Println(entry)
	// Output: <Entry dom="irc" ttl="" cls="IN" typ="A" ["1.2.3.4"]>
	// <Entry dom="irc" ttl="" cls="" typ="A" ["1.2.3.4"]>
	// <Entry dom="irc" ttl="" cls="IN" typ="A" ["1.2.3.4"]>
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
