goats-html: Go Attribute-based Template System for HTML.


What is goats-html?
==================

goats-html is an atribute-based template system specialized for HTML. It borrows
concepts from Template Attribute Language (TAL) (https://en.wikipedia.org/wiki/
Template_Attribute_Language) and implements its own set of attributes. A well-
known TAL system is used in Plone (and Chameleon as a standalone implementation:
http://chameleon.readthedocs.org/en/latest/). Goats-html is largely inspired by
Chameleon.

Differnt from most of other template system, goats-html does not maintain a run-
time rendering engine (AST). Instead, templates are preprocessed and translated
into Go language! For each template there will be a set of Go structs generated.
Your Go program should import and call these templates types hence the template
logic is statically linked into your binary.

Compared to classic template systems like velocity, etc., templates written in
goats-html is more readable and maintainable because of the ATL syntax. However,
due to its not having a runtime rendering engine (AST), it's hard to write
templates in goats-html in dev environment. Even changing one character of the
template forces you to rebuild your server and restart it. For this reason I
introduced a specially designed developer server which enable you to modify your
template without rebuilding/restarting your server. Hence there are dual execution
modes in goats-html: production mode and development mode:

![Dual Modes](images/DualModes.png)

For each template, the command goats generates a Go interface, a template
implementation, and a template proxy, both implements the same interface. When
you compile your server without the flag "--tags gots_devmod" the built binary
contains the template implementation. When it's compiled with the flag, the
stub will convert the template call into a HTTP request to the development
server. So in dev mode you don't need to rebuild/restart your server if you
modified the template (as long as the template interface was not changed).

Both the template generator and the dev server is provided the command goats.


Download the Prebuilt Debian/Ubuntu Package
===========================================

Download the [deb package](https://github.com/linuxerwang/goats-html/blob/master/binaries/goats-html_0.1.0_amd64.deb?raw=true).

Install and Setup
-----------------

    $ sudo dpkg -i goats-html_0.1.0_amd64.deb

Now you need to setup your GOPATH to include the goats sources:

    $ export GOPATH=$GOPATH:/usr/share/goats-html/go

Once it's installed, you can generate GO programs from your template (suppose
your current working directory is in $GOPATH):

    $ goats gen --template_dir goats-html/example

If you are not in $GOPATH, you need to specify the path to the top source
directory:

    $ goats gen --package_root --template_dir goats-html/example

The output directory is by default the same as the template directory, but
you can specify differently:

    $ goats gen --template_dir goats-html/example --output_dir goats-html/mypkg


Build from Source Code
======================

GOATS consists of two parts: the goats command, and the runtime library.
The command "goats" should be put into a place where you can execute without
the absolute path, such as "/usr/bin". The runtime library contains a set of
GO programs which will be compiled into your binary.


Install dependent packages
--------------------------

    $ go get code.google.com/p/go.net/html
    $ go get github.com/howeyc/fsnotify

Build
-----

Checkout the code in your GOPATH, in my case it's ~/go/src:

    $ cd ~/go/src
    $ git clone https://github.com/linuxerwang/goats-html

Setup the GOATS_GOPATH (this is only needed to build GOATS itself, you don't
need this environment for your application):

    $ export GOATS_GOPATH=~/go/src
    $ cd goats-html
    $ make

or to create Debian/Ubuntu packages on x86_64:

    $ make amd64


Run the Example Program
=======================

Under your GOPATH folder (on my machine it's ~/go), run the following commands:

    $ cd $GOPATH
    $ go run goats-html/examples/main.go --benchmark --large

Benchmark example template:

    $ go run goats-html/examples/main.go --benchmark --small
    $ go run goats-html/examples/main.go --benchmark --large


Run the Example Server in Dev Mode
==================================

Under your GOPATH folder (on my machine it's ~/go), run the following commands:

    $ cd $GOPATH
    $ go run --tags goats_devmod goats-html/examples/server/main.go

visit the template: http://localhost:8000.


Tags and Attributes
===================

All goats-html attributes start with "go:". There are many such attributes and an
HTML tag can be attached by multiple attributes. However, not all attributes can
be attached to all HTML tags. If a tag was attached by multiple attributes,
there's an inherent execution order of them, regardless of their attaching order
on the tag. The following section lists all tags and corresponding attributes
with specific execution order.

ANY TAG:

  * go:template
  * go:arg (multiple)
  * go:var (multiple)
  * go:attr (multiple)

&lt;HTML&gt;:

  * go:import

Other TAGs:

  * go:if
  * go:for
  * go:content
  * go:replace
  * go:replaceable
  * go:switch
  * go:case
  * go:default
  * go:call
  * go:omit-tag


Template Built-in Functions:
============================

In goats-html templates, the following built-in functions can be used in
expressions such as go:var, go:arg, go:content:

* center
* cut
* debug
* floatformat
* join
* len
* ljust
* rjust
* title
* quote

We intented to implement filters like in Django template (using | as the filter
operator), and finished a primitive implementation which pass the expression
to go package "ast/parser" and parse the binary operator |. Unfortunately | is
the "bitwise or" in Go lang and it has lower precedence so such expressions
can't be correctly interpreted:

    go:if="price >= 10 && title|length > 20"

because && has higher precedence than |, the express is interpreted as:

    (price >= 10 && title)|(length > 20)

instead of:

    price >= 10 && (title|length) > 20

We don't want to force template authors to use parenthesis, so at present only
built-in function calls are accepted.


Examples
========

Template examples:

    <html>
      <div go:template="ProductCard"
           go:arg="product: proto.Product">
         <div>Name:</div>
         <div go:content="product.Name"></div>
         <div>Price:</div>
         <div go:content="product.Price"></div>
      </div>
    </html>

    <html go:template="HomePage"
          go:import="products/templets.html as product"
          go:arg="pageData: proto.PageData">
    <body>
      <p>My card:
        <div go:template="UserCard"
             go:arg="user: proto.User = pageData.loginUser"
             go:if="user.IsActive"
             go:var="age: time.Now().Year - user.Birthday.Year">
          <span go:content="title(user.Name)"></span>
          <div go:for="@idx, skill: user.skills">
            <span go:content="idx.Counter"><span go:content="skill">
          </div>
          <span>Age:<span> <span go:content="age"></span>
        </div>
        Items I sell:
        <div go:for="product: pageData.products"
             go:call="product#ProductCard">
        </div>
      </p>

      Your Friends:<hr>
      <p go:for="friend: pageData.Friends">
        <span go:content="friend.Name"></span>
        <div go:call="#UserCard"
             go:var="user: friend"></div>
      </p>

      <ul>
        <li go:for="product: pageData.products" go:content="product.Name"></li>
      </ul>

      <div go:for="product: pageData.products">
        <span go:content="product.Name"></span>
        <span go:content="product.Price"></span>
      </div>
    </body>
    </html>


Credits
=======

* linuxerwang (linuxerwang@gmail.com): Created the original system.
* nwlearning (nwlearning@gmail.com): Implemented most of the built-in template functions.


TODO List
=========

* (0.1.1) go:autoescape
* (0.1.x) more built-in functions
*
* (0.2.0) developer server, auto populate and sample populate
* (0.2.0) go:trans
* (0.2.0) support ${} in tag content?
* (0.2.0) template cyclic reference detection
* (0.2.0) template author guide.
* (0.2.0) web framework user guide.
* (0.2.0) csrftoken & other security enhancement?
* (0.2.0) better spaces treatment for tags without go:for.
*
* (LOW PRI) javascript minimize
* (LOW PRI) css minimize
* (LOW PRI) support generic types in template args (GenericArray, GenericMap, GenericStruct)
* (LOW PRI) css class name mangling
* (LOW PRI) properly report go attr misuse
* (LOW PRI) Indent HTML
* (LOW PRI) Ternary operator
* (LOW PRI) customize filters
* (LOW PRI) implement filters with |, need introduce a customized expression parser?

