# Population Projection Model

This file is part of a prototype model for predicting the future population of various types of geographic regions.

The content of this file is focused on a model which allows the prediction of future population over various planning horizons
for a single region.

# Input Data

In order to predict future populations of a specific region we assume the availability of the following data:

## A population growth prediction

A population growth prediction of of the regions population by gender and by age (0 to 90+) as at the dates 1-1-2019 out to the date 1-1-2023. This data can be thought of as
a table of the form:


| Date | Gender | Age | Population |
| -----|--------|-----|------------ |
| 1-1-2018 | F | 0 | xxxx |
| 1-1-2018 | F | 1 | xxxx |
| 1-1-2018 | F | 2 | xxxx |
| 1-1-2018 | F | 3 | xxxx |
| 1-1-2018 | F | .. | xxxx |
| 1-1-2018 | F | 90+ | xxxx |
| 1-1-2018 | M | 0 | xxxx |
| 1-1-2018 | M | 1 | xxxx |
| 1-1-2018 | M | .. | xxxx |
| 1-1-2018 | M | 90+ | xxxx |
| .... | .. | .. | .. |
| 1-1-2023 | F | 0 | xxxx |
| 1-1-2023 | F | 1 | xxxx |
| 1-1-2023 | F | 2 | xxxx |
| 1-1-2023 | F | 3 | xxxx |
| 1-1-2023 | F | .. | xxxx |
| 1-1-2023 | F | 90+ | xxxx |
| 1-1-2023 | M | 0 | xxxx |
| 1-1-2023 | M | 1 | xxxx |
| 1-1-2023 | M | .. | xxxx |
| 1-1-2023 | M | 90+ | xxxx |

## The actual population for a base date

In out case the base date is 1-1-2023. 

In order to perform a populatin prediction into the future we need actual population data by age and gender for this date.

# The goal

The goal is to answer a question such as:

"What will the population of females between the age of 19 and 35 inclusive be as at 1-1-2035"

or

"What will be the population of men and women older than 65 be as at Jan 1st each year for the next 20 years."

# How 
## Growth rates

From the input data we will calculate 3 sequences of growth rates for each age:

- growth rates for each years growth of the female population at each age
- growth rates for each years growth of the male population at each age
- growth rates for each years growth in total population at each age.

That will be 3 * 91 sequences of growth rates.

For each age and each gender there are 6 population values, one for each of 1-1-2018, 1-1-2019, 1-1-2020, 1-1-2021, 1-1-2022, 1-1-2023.

Let the function `population(gender, age)` return a list/slice/array of six population values for each of those dates.

The pseudo code:
```
func growthRate(gender, age) []float {
    g = make([]float64, 5)
    p = population(gender, age)
    for i = 1 to 6 {
        g = append(g, p[i]/p[i-1])
    }
    return g
}
```
calculates the year by year growth rate for that genedeer and that age.

## Forecasts

How to make population forcasts using the growth rates ?

To illustrate the approach lets say we want a 5 and 10 year forecast for a particular gender and age. The following pseudo function will do the heavy lifting.

```
// see model.go for defn of FutureYears
func forcast(gender, age, horizon []FutureYears) [](FutureYear, int) {
    // horizon entries must be in order
    farthest_year = horizon[len(horizon)-1]
    p = population(gender, age)
    pop_now = p[len(p)-1]
    pop_latest = pop_now
    fpop = make([](FutureYear, int))
    starty = 2024
    // apply the growth rates year by year for the first 5 years
    g = growhRate(gender, age)
    for y := 2024; y <= 2028; y++ {
        p = pop_latest * g[y - starty]
        fpop = append((y, p))
    }
    // apply the last growth rate to all years after the 5th
    for y := 2029; y <= farthest_year; y++ {
        p = pop_latest * g[len(g)-1]
        fpop = append((y, p))
    }
    result := make([](FutureYear, int))
    for _,ent := range fpop {
        ytmp = ent[0]
        if slices.BinarySearch(horizon, ytmp) {
            result = append(ent)
        }
    }
    return result
} 

```