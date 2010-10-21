include $(GOROOT)/src/Make.inc

TARG=mamire
GOFILES=mamire.go

all: unlib thread mamirelib mamire

mamire: mamire.$O $(OFILES)
		$(LD) -o $@ mamire.$O $(OFILES)

thread: thread.go
		$(GC) thread.go

unlib: unlib.go
		$(GC) unlib.go

mamirelib: mamire.go unlib.go thread.go
		$(GC) mamire.go

.PHONY: clean
clean:
		rm -rf *.$(O) $(TARG)

