.PHONY: all clean

tmp = *.tex *.pdf
data = privacy-utility-tradeoff.csv
charts = \
	full-vs-partial.tex \
	len-partial-msmt.tex \
	partial-rec-heatmap.tex

all: $(charts)

$(charts): %.tex: plot-%.R
	Rscript $< $(data)

clean:
	rm -f $(tmp)
