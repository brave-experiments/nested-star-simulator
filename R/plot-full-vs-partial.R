require(ggplot2)
require(tikzDevice)

args <- commandArgs(trailingOnly = TRUE)
input_file <- args[1]
data <- read.csv(input_file, header=TRUE)

# Only analyse our partial records.
partial <- subset(data, type == "Partial")

# Turn fraction into percentage.
partial$frac = partial$frac* 100

# Derive full measurements from partial ones.
full <- partial
full$frac <- 100 - full$frac
full$type <- "Full"

data <- rbind(full, partial)

#cairo_pdf("full-vs-partial.pdf", width = 4, height=2.5)
tikz(file = "full-vs-partial.tex", standAlone=F, width = 3.5, height = 2)

ggplot(data, aes(x = k,
                 y = frac,
                 color = type,
                 linetype = type,
                 shape = type)) +
    geom_point(size = 2) +
    geom_line() +
    theme_minimal() +
    labs(x = "k-anonymity threshold",
         y = "\\% of msmts",
         color = "Msmt type",
         linetype = "Msmt type",
         shape = "Msmt type")

dev.off()
