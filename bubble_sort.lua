list = {131,5,344,356,109,179,385,488,36,478,42,438,230,49,271,95,498,410,195,237,453,157,306,493,275,357,150,94,92,57,398,350,458,149,115,60,363,322,392,152,55,98,330,146,229,409,251,105,406,396,258,65,359,182,201,447,200,71,288,332,403,250,35,124,31,23,125,329,391,304,418,295,290,421,394,224,352,89,432,28,1,469,90,461,216,454,314,4,341,297,20,103,72,39,167,221,364,24,233,315,173,265,143,390,122,243,76,283,340,335,227,254,389,18,307,159,428,238,326,255,174,401,449,346,408,91,448,46,445,85,93,45,165,248,123,33,339,497,298,156,334,186,371,62,496,429,235,349,17,107,472,436,348,294,336,132,29,468,328,420,260,483,354,51,338,47,376,300,8,373,415,485,266,380,120,309,463,66,136,97,345,426,145,325,228,102,382,104,213,176,378,284,169,83,164,342,358,16,114,442,246,424,433,466,313,484,368,198,147,194,465,32,25,144,457,19,119,291,381,155,79,82,482,163,106,64,302,462,191,210,324,234,160,203,192,323,88,387,316,37,239,10,141,367,386,267,2,310,183,412,30,116,206,375,166,333,181,87,113,460,481,486,423,193,211,430,175,204,215,74,7,43,437,475,440,117,177,217,489,188,219,413,343,142,78,434,285,139,161,75,388,52,38,162,441,212,417,320,81,225,500,245,370,133,425,68,263,281,12,3,305,492,480,365,450,189,337,209,464,369,419,6,41,151,269,411,73,293,241,110,253,180,207,431,268,187,262,27,470,351,236,48,452,479,223,140,361,379,384,274,196,270,63,327,121,50,153,366,494,287,58,108,446,111,256,400,129,422,393,296,172,9,476,242,226,54,53,158,154,168,353,208,40,264,126,399,257,286,277,303,202,456,171,383,312,362,21,372,499,377,222,487,61,308,244,170,301,214,197,397,100,405,474,59,190,232,471,459,443,278,84,199,67,276,282,69,444,407,185,451,402,137,44,15,311,205,317,467,347,240,22,491,414,148,280,439,321,86,118,395,220,318,80,252,416,259,435,272,404,249,292,360,331,99,14,299,130,273,11,495,374,289,477,13,184,218,96,127,261,247,134,279,56,34,355,138,112,473,135,178,231,26,101,319,128,77,490,70,427,455}
show(list)
print("RUNNING")

function bubbleSort(A)
  local itemCount=#A
  local hasChanged
  repeat
    hasChanged = false
    itemCount=itemCount - 1
    for i = 1, itemCount do
      if A[i] > A[i + 1] then
        A[i], A[i + 1] = A[i + 1], A[i]
        hasChanged = true
				show(A)
      end
    end
  until hasChanged == false
end
bubbleSort(list)
