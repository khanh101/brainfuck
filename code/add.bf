,               Read first character into cell 0
> ,             Read second character into cell 1
<               
-48             Subtract 48 from cell 0 (convert from ASCII)
>               
-48             Subtract 48 from cell 1 (convert from ASCII)
<               
[->+<]          Move cell 0’s value into cell 1 (cell 1 now holds the sum)
>               
+48             Add 48 to cell 1 (convert back to ASCII)
.               Output the resulting character
