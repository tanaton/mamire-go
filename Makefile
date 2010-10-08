include $(GOROOT)/src/Make.inc

TARG=mamire
GOFILES=mamire.go

all: unlib thread $(TARG)

$(TARG): _go_.$O $(OFILES)
		$(LD) -o $@ _go_.$O $(OFILES)

_go_.$O: $(GOFILES) $(PREREQ)
		$(GC) -o $@ $(GOFILES)

thread: thread.go
		$(GC) thread.go

unlib: unlib.go
		$(GC) unlib.go

.PHONY: clean
clean:
		rm -rf _go_.$(O) unlib.$(O) mamire.$(O) thread.$(O) $(TARG)

