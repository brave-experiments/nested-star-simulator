require(ggplot2)
require(tikzDevice)
require(scales)

args <- commandArgs(trailingOnly = TRUE)
input_file <- args[1]
data <- read.csv(input_file, header = TRUE)

part <- subset(data, type == "Partial")
lost <- subset(data, type == "Lost")
full <- subset(data, type == "Full")

# Turn fractions into percentages.
part$pct <- part$frac * 100
lost$pct <- lost$frac * 100
full$pct <- full$frac * 100

cbPalette <- c("#4AA3F1", "#FFC107", "#50ECD2")

data <- rbind(rbind(full, part), lost)

# cairo_pdf("full-vs-partial.pdf", width = 4, height=2.5)
tikz(file = "full-vs-partial.tex", standAlone=F, width = 3.5, height = 2)

ggplot(data, aes(x = k,
                 y = pct,
                 color = type,
                 linetype = type,
                 shape = type)) +
    geom_point(size = 2) +
    geom_line() +
    scale_colour_manual(values=cbPalette) +
    scale_x_continuous(labels=comma, trans="log2") +
    theme_minimal() +
    labs(x = "$K$ (log)",
         y = "\\% of records ",
         color = "Record type",
         linetype = "Record type",
         shape = "Record type")

dev.off()
