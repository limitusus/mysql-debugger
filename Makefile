PROGS=bin/frm-parser

ALL: $(PROGS)

GO := go

bin/frm-parser: src/frm-parser.go
	which go
	$(GO) build -o $@ $<

.PHONY: clean

clean:
	-$(RM) $(PROGS) *~ src/*~
