package gomin

import "testing"

func TestExpected(t *testing.T) {
	checkResult := func(description, subject, expected string) {
		r, err := Js([]byte(subject))

		if err != nil {
			t.Errorf(description+" error: %s", err.Error())
		}

		if r != expected {
			t.Errorf(description+" test failed: %s", r)
		}
	}

	checkResult(
		"Single comment",
		"// This is a comment \n var a = 1",
		"\nvar a=1")

	checkResult(
		"Spacing around vars",
		"var        a       =        1",
		"\nvar a=1")

	checkResult(
		"Multiple single line comments",
		"// This is a comment \n // Another comment \nvar a = 1",
		"\nvar a=1")

	checkResult(
		"Spacing in string preserved",
		"var a = 'its    happening   again'",
		"\nvar a='its    happening   again'")

	checkResult(
		"Vars in comments should be ignored",
		"// This is a comment \n// var b = 2 \n var a = 1",
		"\nvar a=1")

	checkResult(
		"Spacing in functions",
		"func   Bob  (  x  ) {   var   y  =  x; } ",
		"\nfunc Bob(x){var y=x;}")

	checkResult(
		"Spacing in functions (newlines)",
		"func   Bob  (  x  )\n {\n   var   y  =  x;\n }\n ",
		"\nfunc Bob(x)\n{var y=x;}")

	checkResult(
		"Square brackets",
		"var a = b [ c ]; ",
		"\nvar a=b[c];")
	checkResult(
		"Square brackets with newline",
		"var a = b [ c ]\n var d = 123; ",
		"\nvar a=b[c]\nvar d=123;")

	checkResult(
		"Spacing in functions",
		"func   Bob  (  x  ) {   var   y  =  x; } ",
		"\nfunc Bob(x){var y=x;}")

	checkResult(
		"Math division",
		"var a = 3/(v+w)",
		"\nvar a=3/(v+w)")

	checkResult(
		"Comments in strings",
		"var a = 'here // it is'",
		"\nvar a='here // it is'")

	checkResult(
		"Carriage return",
		"var a = 1;\rvar b = 2",
		"\nvar a=1;var b=2")

	checkResult(
		"Carriage return (no semicolon)",
		"var a = 1\rvar b = 2",
		"\nvar a=1\nvar b=2")

	checkResult(
		"Control charcters",
		"var a = 1;\t\bvar b = 2",
		"\nvar a=1;var b=2")
	checkResult(
		"Maths",
		"var a = 1 * 2 + a - f / 7 +(8*8 )",
		"\nvar a=1*2+a-f/7+(8*8)")

	checkResult(
		"Maths",
		"var a = (g++)",
		"\nvar a=(g++)")

	checkResult(
		"Maths",
		"var a = (g++ ++g)",
		"\nvar a=(g++ ++g)")

	checkResult(
		"String literal",
		"var a = 'some\\' thing'",
		"\nvar a='some\\' thing'")

	checkResult(
		"Maths",
		"var a = (g++ \n++g)",
		"\nvar a=(g++\n++g)")

	checkResult(
		"Maths",
		"var a = (g)\n\n",
		"\nvar a=(g)")

	checkResult(
		"Maths",
		"var a = b++//c",
		"\nvar a=b++") //+"+"+"+"+"/"+"c")
	checkResult(
		"Maths",
		"var a = (g "+string('/')+"g)\n\n",
		"\nvar a=(g"+string('/')+"g)")

	checkResult(
		"Maths",
		"var a = 213// a comment'", //+string('/')+string('/')+"g)", //+string('/')+"(g))\n\n",
		"\nvar a=213")

	checkResult(
		"Maths",
		"var a = (213)// a comment'", //+string('/')+string('/')+"g)", //+string('/')+"(g))\n\n",
		"\nvar a=(213)")

	checkResult(
		"Maths",
		"var a = (213+(//() a comment'", //+string('/')+string('/')+"g)", //+string('/')+"(g))\n\n",
		"\nvar a=(213+(")
}

func TestMultiLineComments(t *testing.T) {
	js := " /**                    \n" +
		"    *  This is a comment: \n" +
		"    *  var a = 132        \n" +
		"    */                    \n" +
		"var a = 2"

	r, err := Js([]byte(js))
	if err != nil {
		t.Errorf("Single comment broke: %s", err.Error())
	}
	if r != "\nvar a=2" {
		t.Errorf("Multiline comment failed: %s", r)
	}
}

// This is the example used on the JSMin page.
// See: https://github.com/douglascrockford/JSMin
func TestJSMinExample(t *testing.T) {
	js := "// is.js \n" +
		" \n" +
		"// (c) 2001 Douglas Crockford \n" +
		"// 2001 June 3 \n" +
		" \n" +
		" \n" +
		"// is \n" +
		" \n" +
		"// The -is- object is used to identify the browser.  Every browser edition \n" +
		"// identifies itself, but there is no standard way of doing it, and some of \n" +
		"// the identification is deceptive. This is because the authors of web \n" +
		"// browsers are liars. For example, Microsoft's IE browsers claim to be \n" +
		"// Mozilla 4. Netscape 6 claims to be version 5. \n" +
		" \n" +
		"// Warning: Do not use this awful, awful code or any other thing like it. \n" +
		"// Seriously. \n" +
		" \n" +
		"var is = { \n" +
		"    ie:      navigator.appName == 'Microsoft Internet Explorer', \n" +
		"    java:    navigator.javaEnabled(), \n" +
		"    ns:      navigator.appName == 'Netscape', \n" +
		"    ua:      navigator.userAgent.toLowerCase(), \n" +
		"    version: parseFloat(navigator.appVersion.substr(21)) || \n" +
		"             parseFloat(navigator.appVersion), \n" +
		"    win:     navigator.platform == 'Win32' \n" +
		"} \n" +
		" \n" +
		"is.mac = is.ua.indexOf('mac') >= 0; \n" +
		" \n" +
		"if (is.ua.indexOf('opera') >= 0) { \n" +
		"    is.ie = is.ns = false; \n" +
		"    is.opera = true; \n" +
		"} \n" +
		" \n" +
		"if (is.ua.indexOf('gecko') >= 0) { \n" +
		"    is.ie = is.ns = false; \n" +
		"    is.gecko = true; \n" +
		"} "

	e := "\nvar is={ie:navigator.appName=='Microsoft Internet Explorer',java:navigator.javaEnabled(),ns:navigator.appName=='Netscape',ua:navigator.userAgent.toLowerCase(),version:parseFloat(navigator.appVersion.substr(21))||parseFloat(navigator.appVersion),win:navigator.platform=='Win32'}\n" +
		"is.mac=is.ua.indexOf('mac')>=0;if(is.ua.indexOf('opera')>=0){is.ie=is.ns=false;is.opera=true;}\n" +
		"if(is.ua.indexOf('gecko')>=0){is.ie=is.ns=false;is.gecko=true;}"

	r, err := Js([]byte(js))
	if err != nil {
		t.Errorf("Single comment broke: %s", err.Error())
	}
	if r != e {
		t.Errorf("Original test failed: %s", r)
	}
}
