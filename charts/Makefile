.PHONY: all clean

tmp = *.tex *.pdf *.png
data = privacy-utility-tradeoff.csv
charts = \
	full-vs-partial.tex \
	partial-rec-heatmap.tex \
	unlocked-digits.tex

all: $(charts)

$(charts): %.tex: plot-%.R
	Rscript $< $(data)

clean:
	rm -f $(tmp)
