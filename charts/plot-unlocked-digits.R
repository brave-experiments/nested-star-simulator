require(ggplot2)
require(tikzDevice)
require(scales)

# Our data set contains eight fields: one country code plus seven lat/lon pairs
# of increasing accuracy.
num_fields = 8

args <- commandArgs(trailingOnly = TRUE)
input_file <- args[1]
data <- read.csv(input_file, header = TRUE)

# Initialize two data frames that we are going to populate with our data.
star <- data.frame(k = numeric(),
                   digits = numeric())
nested_star <- data.frame(k = numeric(),
                          digits = numeric())

for (cur_k in unique(data$k)) {
    s <- subset(data, type == "Partial" & k == cur_k)
    partial_frac <- s$frac

    s <- subset(data, type == "LenPartMsmt" & k == cur_k)
    partial_digits <- sum(s$len_part_msmts * s$num_part_msmts * partial_frac)

    s <- subset(data, type == "Full" & k == cur_k)
    mean_digits <- num_fields * s$frac

    nested_star <- rbind(nested_star, list(k=cur_k, digits=mean_digits + partial_digits))
    star <- rbind(star, list(k=cur_k, digits=mean_digits))
}

star$sys <- "$K$-TA schemes"
nested_star$sys <- "\tool"

star$pct <- star$digits / num_fields * 100
nested_star$pct <- nested_star$digits / num_fields * 100

data <- rbind(star, nested_star)

cbPalette <- c("#4AA3F1", "#FFC107", "#50ECD2")

# cairo_pdf("unlocked-digits.pdf", width = 4, height=2.5)
tikz(file = "unlocked-digits.tex", standAlone=F, width = 3.5, height = 2)

sprintf("Largest difference in performance: %.2f%%", max(nested_star$pct - star$pct))

ggplot(data, aes(x = k,
                 y = pct,
                 color = sys,
                 linetype = sys,
                 shape = sys)) +
    geom_point(size = 2) +
    geom_line() +
    scale_colour_manual(values=cbPalette) +
    scale_x_continuous(labels=comma, trans="log2") +
    theme_minimal() +
    labs(x = "k-anonymity threshold (log)",
         y = "\\% of unlocked attributes",
         color = "System",
         linetype = "System",
         shape = "System")

dev.off()
