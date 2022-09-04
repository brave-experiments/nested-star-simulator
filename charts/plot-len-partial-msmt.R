require(ggplot2)
require(tikzDevice)

args <- commandArgs(trailingOnly = TRUE)
input_file <- args[1]
data <- read.csv(input_file, header=TRUE)

# Get rid of partial measurements.
data <- subset(data, type == "LenPartMsmt" & k == 100)

#cairo_pdf("partial-msmt-dist.pdf", width = 4, height=2)
tikz(file = "len-partial-msmt.tex", standAlone=F, width = 3.25, height = 1.75)

ggplot(data, aes(x = len_part_msmts,
                 y = num_part_msmts)) +
    geom_bar(stat="identity") +
    scale_x_continuous(breaks = c(1, 2, 3, 4, 5, 6, 7, 9)) +
    theme_minimal() +
    labs(x = "\\# of unlocked attributes",
         y = "\\# of msmts")

dev.off()
