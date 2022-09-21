.PHONY: all code data charts clean

all: code data charts

code:
	make -C code/

data:
	make -C data/

charts:
	make -C charts/

clean:
	make -C code/ clean
	make -C data/ clean
	make -C charts/ clean
