+++
date = '2025-07-03T06:52:11+02:00'
draft = true
title = 'Mathmetics'
tags = ['computer-science']
+++

Learning Computer Science in my free time, part 1.

Today I wanted to start with OSSU's Computer Science course,
however I noticed it mentioned the pre-requirement of highschool math.
Unfortunately I performed so poorly in my first year after Dutch primary school 
(didn't want to learn french), 
that the highest grade of math I got was pretty sub-par.

So today I made sure to start at the very basics, 
I noticed quickly that most of Arithmics is a bit too basic to study, 
I went over the classes though and picked some out that potentially were unknown to me.
From these classes I made a few notes and jotted them down here; 
perhaps I'll find a better way to visualize these and update this post.
Anyways, probably not interesting to read unless you are me.

> TODO: Visualize these. 

### Calculating the area of triangle
To calculate a area of a triangle, you can take it's base (B) and height (H)
and calculate as follows:

```
B*H/2
10*3/2
15/2
7.5 cm2
```

to calculate a trapezoid
Take the bases's length (paralels) B1, B2 add those together, multiply those by the height H and divide it by 2.

taking a trapezoid with a base of 3cm, other base of 5cm and a height of 2cm (taken from the course)
so 

```
(B1+B2)*H/2
(5 + 3) * 2 / 2
8 * 2 / 2
16 / 2
8 cm2
```

Next something I just forgot, and I would have to look up, calculating the area of a circle (approx, since PI is never fully correct).
example radius of 4m

```
PI * Radius (TODO: to the power of 2, need to see how to type this without copy)
3.14 * 4 power2
3.14 * 16 
50.24 m2
```

Next Sphere volume, I don't remember if this is something I got taught.  

formula: `4 * PI * r3 / 3`

lets assume a radius of 4cm

```
4 * 3.14 * 64 / 3
12.56 * 64 / 3
803.84 / 3
267.94 (and never ending 6) cm3
```

### Calculating circumference

Two formulae can be used, the first for a radius, the second for diameter.
Ofcourse, Diameter is simply multiply Radius by two,
and the other way around is divided by 2.

```
2*PI*R
or
D*PI
```

Example using a diameter of 87 meter, which means the first formula.
We just multiply diameter times PI.

```
87*3.14
273.18 meter
```

Example using 3.6km Radius

```
2*3.6*3.14
7.20*3.14
22.608 km
```

### Volume of a cylinder

Calculating the volume of a cylinder is the same as calculating the surface of a circle;
Then multiply the surface by the height of the cylinder (making it volume instead of surface).

For example a cylinder from 8 * 10 cm where 8 is the diameter of the circle
We need the radius, which is `R=D/2` so it's 4cm

```
PI * r2 * H
3.14 * 16 * 10
50.24 * 10
502.40 cm3
```

### Calculate the volume of a cone

We use the same as the previous section, and divide it by three if it's a cone

```
PI * r2 * H / 3
r.14 * 16 * 10 / 3
50.24 * 10 / 3
502.40 / 3
167.46666666666 cm3

```

### Calculate the value after a percent based decrease

To calculate the original value after a percantage based decrease
current value(`N`) `30` percentage (`R`) `96`. Using the formula below.

```
N:(1-R)
30:(1-0.96)
30:0.04
750
```

6% increase to 25.44

```
N:(1+R)
25.44:(1+0.06)
25.44:1.06

Calc R
    24
106|2544
    254
    ___
   -212
     424
    -424

answer 24
```

### Calculate percent decrease
Calculate percentage decrease from `80` to `30` using the following formula `R=(F-N)/F`.
`R` is the result, `F` is the original value `N` is the decreased value.

```
(F-N)/F
(80-30)/80
50/80
    625
80|50
   500
  -480
   ___
    20
   -16
    __
     400
    -400
     ___
       0

62.5%
```

### Calculate percent increase

Calculate percentage increase from `40` to `70` using the following formula `R=(N-F)/F`.
`R` will be the result, `F` is the original `N` is the increased value.

```
(N-F)/F
(70-40)/70
30/70
70|30
   500
  -480
   ___
    20
   -16
    __
     400
    -400
     ___
       0

62.5%
```

Some notes terminoligy in Statistics/Graphs.
I've used plenty of graph types and I have a good feeling about how to present data;
but the terminology for me was pretty vague, hopefully there's no mistakes here though.

Mean: Actual average of statistics e.g. sum of all items divided by item count.
Median: The absolute middle datapoint, doesn't have to be the mean, when dataset is even, take average of the two middle points
Range: Maximum Value - Minimum value (sum).
Mode: The number that appears most frequent.
Quartile (first and second): take either the first or second half of the data (with uneven, ignore the middle datapoint). Take the Median of this set.
Inter Quartile range:  Second Quartile - Quartile.
A box line can be drawin using the range quartiles and median.

