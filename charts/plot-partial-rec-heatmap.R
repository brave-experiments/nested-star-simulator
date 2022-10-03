require(ggplot2)
require(tikzDevice)
require(scales)

args <- commandArgs(trailingOnly = TRUE)
input_file <- args[1]
data <- read.csv(input_file, header=TRUE)

# Get rid of partial measurements.
data <- subset(data, type == "LenPartMsmt")

# cairo_pdf("partial-rec-heatmap.pdf", width = 3.25, height=2.75)
tikz(file = "partial-rec-heatmap.tex",
     standAlone=F,
     width = 3.25,
     height = 2.75)

# In the Foursquare data set, we have a total of eight location attributes:
# the country code, followed by seven lat/lon pairs of increasing granularity.
numPartialAttrs <- seq(1, 8)
d <- expand.grid(X = numPartialAttrs,
                 Y = unique(data$k))
# Turn fractions into percentages.
d$Z <- data$num_part_msmts * 100

ggplot(data, aes(y = d$Y,
                 x = d$X,
                 fill = d$Z)) +
       geom_tile(color = "light gray") +
       theme_minimal() +
       scale_fill_gradient(low = "white",
                           high = "red") +
       scale_x_discrete(limits = factor(numPartialAttrs)) +
       scale_y_continuous(labels = comma_format(accuracy = 1),
                          trans = "log2",
                          breaks = c(4, 16, 64, 256, 1024, 4096, 16384, 65536)) +
       # Disable the grid.
       theme(panel.grid.major = element_blank(),
             panel.grid.minor = element_blank()) +
       labs(x = "\\# of unlocked location digits",
            y = "$K$ (log)",
            fill = "Percentage")

dev.off()
