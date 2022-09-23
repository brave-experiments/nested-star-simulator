require(ggplot2)
require(tikzDevice)
require(scales)

args <- commandArgs(trailingOnly = TRUE)
input_file <- args[1]
data <- read.csv(input_file, header=TRUE)

# Get rid of partial measurements.
data <- subset(data, type == "LenPartMsmt")

# cairo_pdf("partial-rec-heatmap.pdf", width = 5.5, height=3)
tikz(file = "partial-rec-heatmap.tex",
     standAlone=F,
     width = 3.5,
     height = 3)

# In the Foursquare data set, we have a total of eight location attributes:
# the country code, followed by seven lat/lon pairs of increasing granularity.
numPartialAttrs <- seq(1, 7)
d <- expand.grid(X = numPartialAttrs,
                 Y = unique(data$k))
# Turn fractions into percentages.
d$Z <- data$num_part_msmts * 100

# Add a thousands separator.  Normally, we would use
# scale_x_continuous(labels=comma) but we're dealing with a factor here, so we
# gotta replace the string.
yAxis <- as.factor(d$Y)
levels(yAxis)[levels(yAxis) == "1000"] <- "1,000"

ggplot(data, aes(y = yAxis,
                 x = d$X,
                 fill = d$Z)) +
       geom_tile(color = "gray") +
       theme_minimal() +
       scale_fill_gradient(low = "white", high = "red") +
       scale_x_discrete(limits = factor(numPartialAttrs)) +
       labs(x = "\\# of unlocked location digits",
            y = "k-anonymity threshold",
            fill = "Percentage")

dev.off()
