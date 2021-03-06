//line parser_uri.rl:1
// -*-go-*-
//
// SIP(s) URI parser
package sipmsg

//line parser_uri.rl:7

//line parser_uri.go:12
var _uri_actions []byte = []byte{
	0, 1, 0, 1, 1, 1, 2, 1, 3,
	1, 4, 1, 5, 1, 6, 1, 8,
	1, 9, 2, 0, 6, 2, 4, 5,
	2, 5, 0, 2, 7, 0, 2, 8,
	0, 3, 5, 0, 8, 3, 7, 0,
	8,
}

var _uri_key_offsets []int16 = []int16{
	0, 0, 6, 16, 29, 35, 41, 47,
	53, 59, 65, 72, 80, 88, 96, 98,
	105, 114, 116, 119, 121, 124, 126, 129,
	132, 133, 135, 138, 139, 142, 143, 152,
	161, 169, 177, 185, 193, 195, 201, 210,
	219, 228, 230, 233, 236, 237, 238, 250,
	262, 274, 290, 303, 309, 315, 329, 343,
	349, 355, 362, 370, 377, 385, 391, 398,
	400, 419, 425, 431, 444, 450, 456, 471,
	487, 493, 499, 505, 511, 530, 536, 544,
	550, 558, 564, 572, 580, 588, 596, 604,
	612, 619, 627, 635, 643, 645, 652, 661,
	663, 666, 668, 671, 673, 676, 679, 680,
	683, 684, 687, 688, 697, 706, 714, 722,
	730, 738, 740, 746, 755, 764, 773, 775,
	778, 781, 782, 783, 802, 820, 839, 856,
	874, 888, 912, 918, 924, 930, 936, 953,
	959, 965, 983, 989, 995, 1013, 1031, 1037,
	1043, 1062, 1081, 1087, 1093, 1099, 1105, 1124,
	1130, 1136, 1158, 1175, 1194, 1211, 1230, 1247,
	1266, 1285, 1304, 1323, 1342, 1361, 1371, 1382,
	1394, 1408, 1421, 1437, 1449, 1452, 1456, 1460,
	1464, 1468, 1470, 1481, 1490, 1494, 1498, 1502,
	1506, 1508, 1524, 1539, 1554, 1572, 1590, 1608,
	1626, 1644, 1660, 1681, 1703, 1719, 1741, 1759,
	1777, 1795, 1813, 1831, 1849, 1867, 1885, 1903,
	1921, 1939, 1950, 1961, 1972, 1975, 1994, 2011,
	2027, 2043, 2059, 2075, 2091, 2109, 2128, 2147,
	2166, 2184, 2203, 2222, 2240, 2259, 2279, 2299,
	2319, 2339, 2359, 2377, 2399, 2421, 2443, 2463,
	2483, 2503, 2523, 2543, 2563, 2583, 2603, 2623,
	2643, 2663, 2682, 2701,
}

var _uri_trans_keys []byte = []byte{
	83, 115, 65, 90, 97, 122, 43, 58,
	45, 46, 48, 57, 65, 90, 97, 122,
	33, 37, 47, 61, 93, 95, 126, 36,
	59, 63, 90, 97, 122, 48, 57, 65,
	70, 97, 102, 48, 57, 65, 70, 97,
	102, 48, 57, 65, 70, 97, 102, 48,
	57, 65, 70, 97, 102, 48, 57, 65,
	70, 97, 102, 48, 57, 65, 70, 97,
	102, 58, 48, 57, 65, 70, 97, 102,
	58, 93, 48, 57, 65, 70, 97, 102,
	58, 93, 48, 57, 65, 70, 97, 102,
	58, 93, 48, 57, 65, 70, 97, 102,
	58, 93, 58, 48, 57, 65, 70, 97,
	102, 46, 58, 93, 48, 57, 65, 70,
	97, 102, 48, 57, 46, 48, 57, 48,
	57, 46, 48, 57, 48, 57, 93, 48,
	57, 93, 48, 57, 93, 48, 57, 46,
	48, 57, 46, 46, 48, 57, 46, 46,
	58, 93, 48, 57, 65, 70, 97, 102,
	46, 58, 93, 48, 57, 65, 70, 97,
	102, 58, 93, 48, 57, 65, 70, 97,
	102, 58, 93, 48, 57, 65, 70, 97,
	102, 58, 93, 48, 57, 65, 70, 97,
	102, 58, 93, 48, 57, 65, 70, 97,
	102, 58, 93, 48, 57, 65, 70, 97,
	102, 46, 58, 93, 48, 57, 65, 70,
	97, 102, 46, 58, 93, 48, 57, 65,
	70, 97, 102, 46, 58, 93, 48, 57,
	65, 70, 97, 102, 48, 57, 46, 48,
	57, 46, 48, 57, 46, 58, 43, 58,
	73, 105, 45, 46, 48, 57, 65, 90,
	97, 122, 43, 58, 80, 112, 45, 46,
	48, 57, 65, 90, 97, 122, 43, 58,
	83, 115, 45, 46, 48, 57, 65, 90,
	97, 122, 33, 37, 59, 61, 63, 91,
	95, 126, 36, 47, 48, 57, 65, 90,
	97, 122, 33, 37, 58, 61, 64, 95,
	126, 36, 59, 63, 90, 97, 122, 48,
	57, 65, 70, 97, 102, 48, 57, 65,
	70, 97, 102, 33, 37, 61, 64, 95,
	126, 36, 46, 48, 57, 65, 90, 97,
	122, 33, 37, 61, 64, 95, 126, 36,
	46, 48, 57, 65, 90, 97, 122, 48,
	57, 65, 70, 97, 102, 48, 57, 65,
	70, 97, 102, 91, 48, 57, 65, 90,
	97, 122, 45, 46, 48, 57, 65, 90,
	97, 122, 45, 48, 57, 65, 90, 97,
	122, 45, 46, 48, 57, 65, 90, 97,
	122, 48, 57, 65, 90, 97, 122, 45,
	48, 57, 65, 90, 97, 122, 48, 57,
	33, 37, 77, 84, 85, 93, 95, 109,
	116, 117, 126, 36, 43, 45, 58, 65,
	91, 97, 122, 48, 57, 65, 70, 97,
	102, 48, 57, 65, 70, 97, 102, 33,
	37, 93, 95, 126, 36, 43, 45, 58,
	65, 91, 97, 122, 48, 57, 65, 70,
	97, 102, 48, 57, 65, 70, 97, 102,
	33, 36, 37, 63, 93, 95, 126, 39,
	43, 45, 58, 65, 91, 97, 122, 33,
	36, 37, 61, 63, 93, 95, 126, 39,
	43, 45, 58, 65, 91, 97, 122, 48,
	57, 65, 70, 97, 102, 48, 57, 65,
	70, 97, 102, 48, 57, 65, 70, 97,
	102, 48, 57, 65, 70, 97, 102, 33,
	37, 39, 47, 58, 91, 93, 96, 126,
	36, 41, 42, 43, 45, 57, 65, 90,
	95, 122, 48, 57, 65, 90, 97, 122,
	45, 46, 48, 57, 65, 90, 97, 122,
	48, 57, 65, 90, 97, 122, 45, 46,
	48, 57, 65, 90, 97, 122, 48, 57,
	65, 90, 97, 122, 45, 46, 48, 57,
	65, 90, 97, 122, 45, 46, 48, 57,
	65, 90, 97, 122, 45, 46, 48, 57,
	65, 90, 97, 122, 45, 46, 48, 57,
	65, 90, 97, 122, 45, 46, 48, 57,
	65, 90, 97, 122, 45, 46, 48, 57,
	65, 90, 97, 122, 58, 48, 57, 65,
	70, 97, 102, 58, 93, 48, 57, 65,
	70, 97, 102, 58, 93, 48, 57, 65,
	70, 97, 102, 58, 93, 48, 57, 65,
	70, 97, 102, 58, 93, 58, 48, 57,
	65, 70, 97, 102, 46, 58, 93, 48,
	57, 65, 70, 97, 102, 48, 57, 46,
	48, 57, 48, 57, 46, 48, 57, 48,
	57, 93, 48, 57, 93, 48, 57, 93,
	46, 48, 57, 46, 46, 48, 57, 46,
	46, 58, 93, 48, 57, 65, 70, 97,
	102, 46, 58, 93, 48, 57, 65, 70,
	97, 102, 58, 93, 48, 57, 65, 70,
	97, 102, 58, 93, 48, 57, 65, 70,
	97, 102, 58, 93, 48, 57, 65, 70,
	97, 102, 58, 93, 48, 57, 65, 70,
	97, 102, 58, 93, 48, 57, 65, 70,
	97, 102, 46, 58, 93, 48, 57, 65,
	70, 97, 102, 46, 58, 93, 48, 57,
	65, 70, 97, 102, 46, 58, 93, 48,
	57, 65, 70, 97, 102, 48, 57, 46,
	48, 57, 46, 48, 57, 46, 58, 33,
	37, 45, 46, 58, 59, 61, 63, 64,
	95, 126, 36, 47, 48, 57, 65, 90,
	97, 122, 33, 37, 45, 58, 59, 61,
	63, 64, 95, 126, 36, 47, 48, 57,
	65, 90, 97, 122, 33, 37, 45, 46,
	58, 59, 61, 63, 64, 95, 126, 36,
	47, 48, 57, 65, 90, 97, 122, 33,
	37, 58, 59, 61, 63, 64, 95, 126,
	36, 47, 48, 57, 65, 90, 97, 122,
	33, 37, 45, 58, 59, 61, 63, 64,
	95, 126, 36, 47, 48, 57, 65, 90,
	97, 122, 33, 37, 61, 64, 95, 126,
	36, 46, 48, 57, 65, 90, 97, 122,
	33, 37, 44, 58, 59, 61, 63, 64,
	77, 84, 85, 91, 93, 95, 109, 116,
	117, 126, 36, 57, 65, 90, 97, 122,
	48, 57, 65, 70, 97, 102, 48, 57,
	65, 70, 97, 102, 48, 57, 65, 70,
	97, 102, 48, 57, 65, 70, 97, 102,
	33, 37, 44, 47, 58, 61, 64, 91,
	93, 95, 126, 36, 57, 65, 90, 97,
	122, 48, 57, 65, 70, 97, 102, 48,
	57, 65, 70, 97, 102, 33, 37, 44,
	58, 59, 61, 63, 64, 91, 93, 95,
	126, 36, 57, 65, 90, 97, 122, 48,
	57, 65, 70, 97, 102, 48, 57, 65,
	70, 97, 102, 33, 37, 38, 44, 58,
	59, 61, 64, 91, 93, 95, 126, 36,
	57, 63, 90, 97, 122, 33, 37, 38,
	44, 58, 59, 61, 64, 91, 93, 95,
	126, 36, 57, 63, 90, 97, 122, 48,
	57, 65, 70, 97, 102, 48, 57, 65,
	70, 97, 102, 33, 37, 38, 44, 47,
	58, 61, 63, 64, 91, 93, 95, 126,
	36, 57, 65, 90, 97, 122, 33, 37,
	38, 44, 47, 58, 61, 63, 64, 91,
	93, 95, 126, 36, 57, 65, 90, 97,
	122, 48, 57, 65, 70, 97, 102, 48,
	57, 65, 70, 97, 102, 48, 57, 65,
	70, 97, 102, 48, 57, 65, 70, 97,
	102, 33, 37, 38, 44, 47, 58, 61,
	63, 64, 91, 93, 95, 126, 36, 57,
	65, 90, 97, 122, 48, 57, 65, 70,
	97, 102, 48, 57, 65, 70, 97, 102,
	33, 37, 39, 44, 47, 58, 59, 61,
	63, 64, 91, 93, 96, 126, 36, 41,
	42, 57, 65, 90, 95, 122, 33, 37,
	58, 59, 61, 63, 64, 95, 126, 36,
	47, 48, 57, 65, 90, 97, 122, 33,
	37, 45, 46, 58, 59, 61, 63, 64,
	95, 126, 36, 47, 48, 57, 65, 90,
	97, 122, 33, 37, 58, 59, 61, 63,
	64, 95, 126, 36, 47, 48, 57, 65,
	90, 97, 122, 33, 37, 45, 46, 58,
	59, 61, 63, 64, 95, 126, 36, 47,
	48, 57, 65, 90, 97, 122, 33, 37,
	58, 59, 61, 63, 64, 95, 126, 36,
	47, 48, 57, 65, 90, 97, 122, 33,
	37, 45, 46, 58, 59, 61, 63, 64,
	95, 126, 36, 47, 48, 57, 65, 90,
	97, 122, 33, 37, 45, 46, 58, 59,
	61, 63, 64, 95, 126, 36, 47, 48,
	57, 65, 90, 97, 122, 33, 37, 45,
	46, 58, 59, 61, 63, 64, 95, 126,
	36, 47, 48, 57, 65, 90, 97, 122,
	33, 37, 45, 46, 58, 59, 61, 63,
	64, 95, 126, 36, 47, 48, 57, 65,
	90, 97, 122, 33, 37, 45, 46, 58,
	59, 61, 63, 64, 95, 126, 36, 47,
	48, 57, 65, 90, 97, 122, 33, 37,
	45, 46, 58, 59, 61, 63, 64, 95,
	126, 36, 47, 48, 57, 65, 90, 97,
	122, 43, 58, 45, 46, 48, 57, 65,
	90, 97, 122, 33, 37, 61, 95, 126,
	36, 59, 63, 90, 97, 122, 33, 37,
	47, 61, 95, 126, 36, 59, 63, 90,
	97, 122, 33, 37, 58, 61, 64, 91,
	95, 126, 36, 59, 63, 90, 97, 122,
	33, 37, 58, 61, 64, 95, 126, 36,
	59, 63, 90, 97, 122, 33, 37, 47,
	61, 63, 64, 95, 126, 36, 57, 58,
	59, 65, 90, 97, 122, 33, 37, 61,
	91, 95, 126, 36, 59, 63, 90, 97,
	122, 47, 58, 63, 47, 63, 48, 57,
	47, 63, 48, 57, 47, 63, 48, 57,
	47, 63, 48, 57, 47, 63, 45, 46,
	58, 59, 63, 48, 57, 65, 90, 97,
	122, 58, 59, 63, 48, 57, 65, 90,
	97, 122, 59, 63, 48, 57, 59, 63,
	48, 57, 59, 63, 48, 57, 59, 63,
	48, 57, 59, 63, 33, 37, 59, 61,
	63, 93, 95, 126, 36, 43, 45, 58,
	65, 91, 97, 122, 33, 37, 59, 63,
	93, 95, 126, 36, 43, 45, 58, 65,
	91, 97, 122, 33, 37, 38, 63, 93,
	95, 126, 36, 43, 45, 58, 65, 91,
	97, 122, 33, 37, 59, 61, 63, 69,
	93, 95, 101, 126, 36, 43, 45, 58,
	65, 91, 97, 122, 33, 37, 59, 61,
	63, 84, 93, 95, 116, 126, 36, 43,
	45, 58, 65, 91, 97, 122, 33, 37,
	59, 61, 63, 72, 93, 95, 104, 126,
	36, 43, 45, 58, 65, 91, 97, 122,
	33, 37, 59, 61, 63, 79, 93, 95,
	111, 126, 36, 43, 45, 58, 65, 91,
	97, 122, 33, 37, 59, 61, 63, 68,
	93, 95, 100, 126, 36, 43, 45, 58,
	65, 91, 97, 122, 33, 37, 59, 61,
	63, 93, 95, 126, 36, 43, 45, 58,
	65, 91, 97, 122, 33, 37, 39, 47,
	58, 59, 63, 91, 93, 96, 126, 36,
	41, 42, 43, 45, 57, 65, 90, 95,
	122, 33, 37, 39, 59, 63, 126, 42,
	43, 45, 46, 48, 57, 65, 70, 71,
	90, 95, 96, 97, 102, 103, 122, 33,
	37, 39, 59, 63, 126, 42, 43, 45,
	46, 48, 57, 65, 90, 95, 122, 33,
	37, 39, 59, 63, 126, 42, 43, 45,
	46, 48, 57, 65, 70, 71, 90, 95,
	96, 97, 102, 103, 122, 33, 37, 59,
	61, 63, 82, 93, 95, 114, 126, 36,
	43, 45, 58, 65, 91, 97, 122, 33,
	37, 59, 61, 63, 65, 93, 95, 97,
	126, 36, 43, 45, 58, 66, 91, 98,
	122, 33, 37, 59, 61, 63, 78, 93,
	95, 110, 126, 36, 43, 45, 58, 65,
	91, 97, 122, 33, 37, 59, 61, 63,
	83, 93, 95, 115, 126, 36, 43, 45,
	58, 65, 91, 97, 122, 33, 37, 59,
	61, 63, 80, 93, 95, 112, 126, 36,
	43, 45, 58, 65, 91, 97, 122, 33,
	37, 59, 61, 63, 79, 93, 95, 111,
	126, 36, 43, 45, 58, 65, 91, 97,
	122, 33, 37, 59, 61, 63, 82, 93,
	95, 114, 126, 36, 43, 45, 58, 65,
	91, 97, 122, 33, 37, 59, 61, 63,
	84, 93, 95, 116, 126, 36, 43, 45,
	58, 65, 91, 97, 122, 33, 37, 59,
	61, 63, 83, 93, 95, 115, 126, 36,
	43, 45, 58, 65, 91, 97, 122, 33,
	37, 59, 61, 63, 69, 93, 95, 101,
	126, 36, 43, 45, 58, 65, 91, 97,
	122, 33, 37, 59, 61, 63, 82, 93,
	95, 114, 126, 36, 43, 45, 58, 65,
	91, 97, 122, 45, 46, 58, 59, 63,
	48, 57, 65, 90, 97, 122, 45, 46,
	58, 59, 63, 48, 57, 65, 90, 97,
	122, 45, 46, 58, 59, 63, 48, 57,
	65, 90, 97, 122, 58, 59, 63, 33,
	37, 45, 46, 58, 59, 61, 63, 64,
	95, 126, 36, 47, 48, 57, 65, 90,
	97, 122, 33, 37, 58, 59, 61, 63,
	64, 95, 126, 36, 47, 48, 57, 65,
	90, 97, 122, 33, 37, 59, 61, 63,
	64, 95, 126, 36, 46, 48, 57, 65,
	90, 97, 122, 33, 37, 59, 61, 63,
	64, 95, 126, 36, 46, 48, 57, 65,
	90, 97, 122, 33, 37, 59, 61, 63,
	64, 95, 126, 36, 46, 48, 57, 65,
	90, 97, 122, 33, 37, 59, 61, 63,
	64, 95, 126, 36, 46, 48, 57, 65,
	90, 97, 122, 33, 37, 59, 61, 63,
	64, 95, 126, 36, 46, 48, 57, 65,
	90, 97, 122, 33, 37, 44, 58, 59,
	61, 63, 64, 91, 93, 95, 126, 36,
	57, 65, 90, 97, 122, 33, 37, 44,
	47, 58, 59, 61, 63, 64, 91, 93,
	95, 126, 36, 57, 65, 90, 97, 122,
	33, 37, 44, 47, 58, 59, 61, 63,
	64, 91, 93, 95, 126, 36, 57, 65,
	90, 97, 122, 33, 37, 44, 47, 58,
	59, 61, 63, 64, 91, 93, 95, 126,
	36, 57, 65, 90, 97, 122, 33, 37,
	44, 58, 59, 61, 63, 64, 91, 93,
	95, 126, 36, 57, 65, 90, 97, 122,
	33, 37, 44, 47, 58, 59, 61, 63,
	64, 91, 93, 95, 126, 36, 57, 65,
	90, 97, 122, 33, 37, 38, 44, 47,
	58, 61, 63, 64, 91, 93, 95, 126,
	36, 57, 65, 90, 97, 122, 33, 37,
	38, 44, 58, 59, 61, 64, 91, 93,
	95, 126, 36, 57, 63, 90, 97, 122,
	33, 37, 38, 44, 47, 58, 61, 63,
	64, 91, 93, 95, 126, 36, 57, 65,
	90, 97, 122, 33, 37, 44, 58, 59,
	61, 63, 64, 69, 91, 93, 95, 101,
	126, 36, 57, 65, 90, 97, 122, 33,
	37, 44, 58, 59, 61, 63, 64, 84,
	91, 93, 95, 116, 126, 36, 57, 65,
	90, 97, 122, 33, 37, 44, 58, 59,
	61, 63, 64, 72, 91, 93, 95, 104,
	126, 36, 57, 65, 90, 97, 122, 33,
	37, 44, 58, 59, 61, 63, 64, 79,
	91, 93, 95, 111, 126, 36, 57, 65,
	90, 97, 122, 33, 37, 44, 58, 59,
	61, 63, 64, 68, 91, 93, 95, 100,
	126, 36, 57, 65, 90, 97, 122, 33,
	37, 44, 58, 59, 61, 63, 64, 91,
	93, 95, 126, 36, 57, 65, 90, 97,
	122, 33, 37, 39, 44, 47, 58, 59,
	61, 63, 64, 91, 93, 96, 126, 36,
	41, 42, 57, 65, 90, 95, 122, 33,
	37, 39, 59, 63, 126, 42, 43, 45,
	46, 48, 57, 65, 70, 71, 90, 95,
	96, 97, 102, 103, 122, 33, 37, 39,
	59, 63, 126, 42, 43, 45, 46, 48,
	57, 65, 70, 71, 90, 95, 96, 97,
	102, 103, 122, 33, 37, 44, 58, 59,
	61, 63, 64, 82, 91, 93, 95, 114,
	126, 36, 57, 65, 90, 97, 122, 33,
	37, 44, 58, 59, 61, 63, 64, 65,
	91, 93, 95, 97, 126, 36, 57, 66,
	90, 98, 122, 33, 37, 44, 58, 59,
	61, 63, 64, 78, 91, 93, 95, 110,
	126, 36, 57, 65, 90, 97, 122, 33,
	37, 44, 58, 59, 61, 63, 64, 83,
	91, 93, 95, 115, 126, 36, 57, 65,
	90, 97, 122, 33, 37, 44, 58, 59,
	61, 63, 64, 80, 91, 93, 95, 112,
	126, 36, 57, 65, 90, 97, 122, 33,
	37, 44, 58, 59, 61, 63, 64, 79,
	91, 93, 95, 111, 126, 36, 57, 65,
	90, 97, 122, 33, 37, 44, 58, 59,
	61, 63, 64, 82, 91, 93, 95, 114,
	126, 36, 57, 65, 90, 97, 122, 33,
	37, 44, 58, 59, 61, 63, 64, 84,
	91, 93, 95, 116, 126, 36, 57, 65,
	90, 97, 122, 33, 37, 44, 58, 59,
	61, 63, 64, 83, 91, 93, 95, 115,
	126, 36, 57, 65, 90, 97, 122, 33,
	37, 44, 58, 59, 61, 63, 64, 69,
	91, 93, 95, 101, 126, 36, 57, 65,
	90, 97, 122, 33, 37, 44, 58, 59,
	61, 63, 64, 82, 91, 93, 95, 114,
	126, 36, 57, 65, 90, 97, 122, 33,
	37, 45, 46, 58, 59, 61, 63, 64,
	95, 126, 36, 47, 48, 57, 65, 90,
	97, 122, 33, 37, 45, 46, 58, 59,
	61, 63, 64, 95, 126, 36, 47, 48,
	57, 65, 90, 97, 122, 33, 37, 45,
	46, 58, 59, 61, 63, 64, 95, 126,
	36, 47, 48, 57, 65, 90, 97, 122,
}

var _uri_single_lengths []byte = []byte{
	0, 2, 2, 7, 0, 0, 0, 0,
	0, 0, 1, 2, 2, 2, 2, 1,
	3, 0, 1, 0, 1, 0, 1, 1,
	1, 0, 1, 1, 1, 1, 3, 3,
	2, 2, 2, 2, 2, 0, 3, 3,
	3, 0, 1, 1, 1, 1, 4, 4,
	4, 8, 7, 0, 0, 6, 6, 0,
	0, 1, 2, 1, 2, 0, 1, 0,
	11, 0, 0, 5, 0, 0, 7, 8,
	0, 0, 0, 0, 9, 0, 2, 0,
	2, 0, 2, 2, 2, 2, 2, 2,
	1, 2, 2, 2, 2, 1, 3, 0,
	1, 0, 1, 0, 1, 1, 1, 1,
	1, 1, 1, 3, 3, 2, 2, 2,
	2, 2, 0, 3, 3, 3, 0, 1,
	1, 1, 1, 11, 10, 11, 9, 10,
	6, 18, 0, 0, 0, 0, 11, 0,
	0, 12, 0, 0, 12, 12, 0, 0,
	13, 13, 0, 0, 0, 0, 13, 0,
	0, 14, 9, 11, 9, 11, 9, 11,
	11, 11, 11, 11, 11, 2, 5, 6,
	8, 7, 8, 6, 3, 2, 2, 2,
	2, 2, 5, 3, 2, 2, 2, 2,
	2, 8, 7, 7, 10, 10, 10, 10,
	10, 8, 11, 6, 6, 6, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 5, 5, 5, 3, 11, 9, 8,
	8, 8, 8, 8, 12, 13, 13, 13,
	12, 13, 13, 12, 13, 14, 14, 14,
	14, 14, 12, 14, 6, 6, 14, 14,
	14, 14, 14, 14, 14, 14, 14, 14,
	14, 11, 11, 11,
}

var _uri_range_lengths []byte = []byte{
	0, 2, 4, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 0, 3,
	3, 1, 1, 1, 1, 1, 1, 1,
	0, 1, 1, 0, 1, 0, 3, 3,
	3, 3, 3, 3, 0, 3, 3, 3,
	3, 1, 1, 1, 0, 0, 4, 4,
	4, 4, 3, 3, 3, 4, 4, 3,
	3, 3, 3, 3, 3, 3, 3, 1,
	4, 3, 3, 4, 3, 3, 4, 4,
	3, 3, 3, 3, 5, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 0, 3, 3, 1,
	1, 1, 1, 1, 1, 1, 0, 1,
	0, 1, 0, 3, 3, 3, 3, 3,
	3, 0, 3, 3, 3, 3, 1, 1,
	1, 0, 0, 4, 4, 4, 4, 4,
	4, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3,
	3, 4, 4, 4, 4, 4, 4, 4,
	4, 4, 4, 4, 4, 4, 3, 3,
	3, 3, 4, 3, 0, 1, 1, 1,
	1, 0, 3, 3, 1, 1, 1, 1,
	0, 4, 4, 4, 4, 4, 4, 4,
	4, 4, 5, 8, 5, 8, 4, 4,
	4, 4, 4, 4, 4, 4, 4, 4,
	4, 3, 3, 3, 0, 4, 4, 4,
	4, 4, 4, 4, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 4, 8, 8, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3,
	3, 4, 4, 4,
}

var _uri_index_offsets []int16 = []int16{
	0, 0, 5, 12, 23, 27, 31, 35,
	39, 43, 47, 52, 58, 64, 70, 73,
	78, 85, 87, 90, 92, 95, 97, 100,
	103, 105, 107, 110, 112, 115, 117, 124,
	131, 137, 143, 149, 155, 158, 162, 169,
	176, 183, 185, 188, 191, 193, 195, 204,
	213, 222, 235, 246, 250, 254, 265, 276,
	280, 284, 289, 295, 300, 306, 310, 315,
	317, 333, 337, 341, 351, 355, 359, 371,
	384, 388, 392, 396, 400, 415, 419, 425,
	429, 435, 439, 445, 451, 457, 463, 469,
	475, 480, 486, 492, 498, 501, 506, 513,
	515, 518, 520, 523, 525, 528, 531, 533,
	536, 538, 541, 543, 550, 557, 563, 569,
	575, 581, 584, 588, 595, 602, 609, 611,
	614, 617, 619, 621, 637, 652, 668, 682,
	697, 708, 730, 734, 738, 742, 746, 761,
	765, 769, 785, 789, 793, 809, 825, 829,
	833, 850, 867, 871, 875, 879, 883, 900,
	904, 908, 927, 941, 957, 971, 987, 1001,
	1017, 1033, 1049, 1065, 1081, 1097, 1104, 1113,
	1123, 1135, 1146, 1159, 1169, 1173, 1177, 1181,
	1185, 1189, 1192, 1201, 1208, 1212, 1216, 1220,
	1224, 1227, 1240, 1252, 1264, 1279, 1294, 1309,
	1324, 1339, 1352, 1369, 1384, 1396, 1411, 1426,
	1441, 1456, 1471, 1486, 1501, 1516, 1531, 1546,
	1561, 1576, 1585, 1594, 1603, 1607, 1623, 1637,
	1650, 1663, 1676, 1689, 1702, 1718, 1735, 1752,
	1769, 1785, 1802, 1819, 1835, 1852, 1870, 1888,
	1906, 1924, 1942, 1958, 1977, 1992, 2007, 2025,
	2043, 2061, 2079, 2097, 2115, 2133, 2151, 2169,
	2187, 2205, 2221, 2237,
}

var _uri_trans_targs []byte = []byte{
	46, 46, 2, 2, 0, 2, 3, 2,
	2, 2, 2, 0, 166, 4, 167, 166,
	166, 166, 166, 166, 166, 166, 0, 5,
	5, 5, 0, 166, 166, 166, 0, 7,
	7, 7, 0, 169, 169, 169, 0, 9,
	9, 9, 0, 170, 170, 170, 0, 45,
	11, 11, 11, 0, 15, 172, 12, 12,
	12, 0, 15, 172, 13, 13, 13, 0,
	15, 172, 14, 14, 14, 0, 15, 172,
	0, 32, 16, 11, 11, 0, 17, 15,
	172, 30, 12, 12, 0, 18, 0, 19,
	28, 0, 20, 0, 21, 26, 0, 22,
	0, 172, 23, 0, 172, 24, 0, 172,
	0, 173, 0, 21, 27, 0, 21, 0,
	19, 29, 0, 19, 0, 17, 15, 172,
	31, 13, 13, 0, 17, 15, 172, 14,
	14, 14, 0, 41, 172, 33, 33, 33,
	0, 37, 172, 34, 34, 34, 0, 37,
	172, 35, 35, 35, 0, 37, 172, 36,
	36, 36, 0, 37, 172, 0, 38, 33,
	33, 0, 17, 37, 172, 39, 34, 34,
	0, 17, 37, 172, 40, 35, 35, 0,
	17, 37, 172, 36, 36, 36, 0, 42,
	0, 17, 43, 0, 17, 44, 0, 17,
	0, 32, 0, 2, 3, 47, 47, 2,
	2, 2, 2, 0, 2, 3, 48, 48,
	2, 2, 2, 2, 0, 2, 49, 165,
	165, 2, 2, 2, 2, 0, 50, 51,
	50, 50, 50, 88, 50, 50, 50, 123,
	213, 213, 0, 50, 51, 53, 50, 57,
	50, 50, 50, 50, 50, 0, 52, 52,
	52, 0, 50, 50, 50, 0, 54, 55,
	54, 57, 54, 54, 54, 54, 54, 54,
	0, 54, 55, 54, 57, 54, 54, 54,
	54, 54, 54, 0, 56, 56, 56, 0,
	54, 54, 54, 0, 88, 58, 178, 178,
	0, 59, 77, 86, 60, 60, 0, 59,
	60, 60, 60, 0, 59, 61, 60, 60,
	60, 0, 60, 178, 178, 0, 62, 178,
	178, 178, 0, 180, 0, 185, 65, 188,
	198, 206, 185, 185, 188, 198, 206, 185,
	185, 185, 185, 185, 0, 66, 66, 66,
	0, 185, 185, 185, 0, 186, 68, 186,
	186, 186, 186, 186, 186, 186, 0, 69,
	69, 69, 0, 186, 186, 186, 0, 71,
	71, 72, 71, 71, 71, 71, 71, 71,
	71, 71, 0, 71, 71, 72, 187, 71,
	71, 71, 71, 71, 71, 71, 71, 0,
	73, 73, 73, 0, 71, 71, 71, 0,
	75, 75, 75, 0, 187, 187, 187, 0,
	194, 195, 194, 186, 186, 186, 186, 196,
	194, 186, 194, 194, 194, 194, 0, 78,
	178, 178, 0, 59, 79, 84, 60, 60,
	0, 80, 178, 178, 0, 59, 81, 82,
	60, 60, 0, 209, 178, 178, 0, 59,
	81, 83, 60, 60, 0, 59, 81, 60,
	60, 60, 0, 59, 79, 85, 60, 60,
	0, 59, 79, 60, 60, 60, 0, 59,
	77, 87, 60, 60, 0, 59, 77, 60,
	60, 60, 0, 122, 89, 89, 89, 0,
	93, 212, 90, 90, 90, 0, 93, 212,
	91, 91, 91, 0, 93, 212, 92, 92,
	92, 0, 93, 212, 0, 109, 94, 89,
	89, 0, 95, 93, 212, 107, 90, 90,
	0, 96, 0, 97, 105, 0, 98, 0,
	99, 103, 0, 100, 0, 212, 101, 0,
	212, 102, 0, 212, 0, 99, 104, 0,
	99, 0, 97, 106, 0, 97, 0, 95,
	93, 212, 108, 91, 91, 0, 95, 93,
	212, 92, 92, 92, 0, 118, 212, 110,
	110, 110, 0, 114, 212, 111, 111, 111,
	0, 114, 212, 112, 112, 112, 0, 114,
	212, 113, 113, 113, 0, 114, 212, 0,
	115, 110, 110, 0, 95, 114, 212, 116,
	111, 111, 0, 95, 114, 212, 117, 112,
	112, 0, 95, 114, 212, 113, 113, 113,
	0, 119, 0, 95, 120, 0, 95, 121,
	0, 95, 0, 109, 0, 50, 51, 124,
	154, 53, 50, 50, 50, 57, 50, 50,
	50, 163, 125, 125, 0, 50, 51, 124,
	53, 50, 50, 50, 57, 50, 50, 50,
	125, 125, 125, 0, 50, 51, 124, 126,
	53, 50, 50, 50, 57, 50, 50, 50,
	125, 125, 125, 0, 50, 51, 53, 50,
	50, 50, 57, 50, 50, 50, 125, 213,
	213, 0, 50, 51, 127, 53, 50, 50,
	50, 57, 50, 50, 50, 213, 213, 213,
	0, 54, 55, 54, 57, 54, 54, 54,
	215, 54, 54, 0, 220, 130, 50, 221,
	50, 50, 50, 57, 229, 238, 246, 185,
	185, 220, 229, 238, 246, 220, 220, 220,
	220, 0, 131, 131, 131, 0, 220, 220,
	220, 0, 133, 133, 133, 0, 222, 222,
	222, 0, 223, 135, 54, 186, 186, 54,
	57, 186, 186, 223, 223, 223, 223, 223,
	0, 136, 136, 136, 0, 223, 223, 223,
	0, 224, 138, 50, 225, 50, 50, 50,
	57, 186, 186, 224, 224, 224, 224, 224,
	0, 139, 139, 139, 0, 224, 224, 224,
	0, 141, 142, 50, 50, 144, 50, 50,
	57, 71, 71, 141, 141, 141, 141, 141,
	0, 141, 142, 50, 50, 144, 50, 227,
	57, 71, 71, 141, 141, 141, 141, 141,
	0, 143, 143, 143, 0, 141, 141, 141,
	0, 145, 146, 54, 54, 71, 71, 226,
	71, 57, 71, 71, 145, 145, 145, 145,
	145, 0, 145, 146, 54, 54, 71, 71,
	226, 71, 57, 71, 71, 145, 145, 145,
	145, 145, 0, 147, 147, 147, 0, 145,
	145, 145, 0, 149, 149, 149, 0, 226,
	226, 226, 0, 145, 146, 54, 54, 71,
	71, 54, 71, 57, 71, 71, 145, 145,
	145, 145, 145, 0, 152, 152, 152, 0,
	227, 227, 227, 0, 235, 236, 235, 50,
	224, 225, 50, 50, 50, 57, 186, 186,
	196, 235, 224, 235, 235, 235, 0, 50,
	51, 53, 50, 50, 50, 57, 50, 50,
	50, 155, 213, 213, 0, 50, 51, 124,
	156, 53, 50, 50, 50, 57, 50, 50,
	50, 161, 125, 125, 0, 50, 51, 53,
	50, 50, 50, 57, 50, 50, 50, 157,
	213, 213, 0, 50, 51, 124, 158, 53,
	50, 50, 50, 57, 50, 50, 50, 159,
	125, 125, 0, 50, 51, 53, 50, 50,
	50, 57, 50, 50, 50, 249, 213, 213,
	0, 50, 51, 124, 158, 53, 50, 50,
	50, 57, 50, 50, 50, 160, 125, 125,
	0, 50, 51, 124, 158, 53, 50, 50,
	50, 57, 50, 50, 50, 125, 125, 125,
	0, 50, 51, 124, 156, 53, 50, 50,
	50, 57, 50, 50, 50, 162, 125, 125,
	0, 50, 51, 124, 156, 53, 50, 50,
	50, 57, 50, 50, 50, 125, 125, 125,
	0, 50, 51, 124, 154, 53, 50, 50,
	50, 57, 50, 50, 50, 164, 125, 125,
	0, 50, 51, 124, 154, 53, 50, 50,
	50, 57, 50, 50, 50, 125, 125, 125,
	0, 2, 49, 2, 2, 2, 2, 0,
	166, 4, 166, 166, 166, 166, 166, 166,
	0, 166, 4, 168, 166, 166, 166, 166,
	166, 166, 0, 169, 6, 166, 169, 166,
	10, 169, 169, 169, 169, 169, 0, 169,
	6, 170, 169, 171, 169, 169, 169, 169,
	169, 0, 170, 8, 166, 170, 166, 171,
	170, 170, 170, 166, 170, 170, 0, 166,
	4, 166, 10, 166, 166, 166, 166, 166,
	0, 166, 25, 166, 0, 166, 166, 174,
	0, 166, 166, 175, 0, 166, 166, 176,
	0, 166, 166, 177, 0, 166, 166, 0,
	62, 179, 63, 64, 70, 178, 178, 178,
	0, 63, 64, 70, 60, 178, 178, 0,
	64, 70, 181, 0, 64, 70, 182, 0,
	64, 70, 183, 0, 64, 70, 184, 0,
	64, 70, 0, 185, 65, 64, 67, 70,
	185, 185, 185, 185, 185, 185, 185, 0,
	186, 68, 64, 70, 186, 186, 186, 186,
	186, 186, 186, 0, 187, 74, 70, 187,
	187, 187, 187, 187, 187, 187, 187, 0,
	185, 65, 64, 67, 70, 189, 185, 185,
	189, 185, 185, 185, 185, 185, 0, 185,
	65, 64, 67, 70, 190, 185, 185, 190,
	185, 185, 185, 185, 185, 0, 185, 65,
	64, 67, 70, 191, 185, 185, 191, 185,
	185, 185, 185, 185, 0, 185, 65, 64,
	67, 70, 192, 185, 185, 192, 185, 185,
	185, 185, 185, 0, 185, 65, 64, 67,
	70, 193, 185, 185, 193, 185, 185, 185,
	185, 185, 0, 185, 65, 64, 76, 70,
	185, 185, 185, 185, 185, 185, 185, 0,
	194, 195, 194, 186, 186, 64, 70, 186,
	186, 196, 194, 186, 194, 194, 194, 194,
	0, 196, 196, 196, 64, 70, 196, 196,
	196, 197, 197, 196, 196, 197, 196, 0,
	196, 196, 196, 64, 70, 196, 196, 196,
	196, 196, 196, 0, 196, 196, 196, 64,
	70, 196, 196, 196, 194, 194, 196, 196,
	194, 196, 0, 185, 65, 64, 67, 70,
	199, 185, 185, 199, 185, 185, 185, 185,
	185, 0, 185, 65, 64, 67, 70, 200,
	185, 185, 200, 185, 185, 185, 185, 185,
	0, 185, 65, 64, 67, 70, 201, 185,
	185, 201, 185, 185, 185, 185, 185, 0,
	185, 65, 64, 67, 70, 202, 185, 185,
	202, 185, 185, 185, 185, 185, 0, 185,
	65, 64, 67, 70, 203, 185, 185, 203,
	185, 185, 185, 185, 185, 0, 185, 65,
	64, 67, 70, 204, 185, 185, 204, 185,
	185, 185, 185, 185, 0, 185, 65, 64,
	67, 70, 205, 185, 185, 205, 185, 185,
	185, 185, 185, 0, 185, 65, 64, 67,
	70, 193, 185, 185, 193, 185, 185, 185,
	185, 185, 0, 185, 65, 64, 67, 70,
	207, 185, 185, 207, 185, 185, 185, 185,
	185, 0, 185, 65, 64, 67, 70, 208,
	185, 185, 208, 185, 185, 185, 185, 185,
	0, 185, 65, 64, 67, 70, 193, 185,
	185, 193, 185, 185, 185, 185, 185, 0,
	59, 61, 63, 64, 70, 210, 60, 60,
	0, 59, 61, 63, 64, 70, 211, 60,
	60, 0, 59, 61, 63, 64, 70, 60,
	60, 60, 0, 63, 64, 70, 0, 50,
	51, 127, 214, 128, 129, 50, 140, 57,
	50, 50, 50, 213, 213, 213, 0, 50,
	51, 128, 129, 50, 140, 57, 50, 50,
	50, 125, 213, 213, 0, 54, 55, 64,
	54, 70, 57, 54, 54, 54, 216, 54,
	54, 0, 54, 55, 64, 54, 70, 57,
	54, 54, 54, 217, 54, 54, 0, 54,
	55, 64, 54, 70, 57, 54, 54, 54,
	218, 54, 54, 0, 54, 55, 64, 54,
	70, 57, 54, 54, 54, 219, 54, 54,
	0, 54, 55, 64, 54, 70, 57, 54,
	54, 54, 54, 54, 54, 0, 220, 130,
	50, 221, 129, 137, 140, 57, 185, 185,
	220, 220, 220, 220, 220, 0, 222, 132,
	54, 185, 185, 64, 134, 70, 57, 185,
	185, 222, 222, 222, 222, 222, 0, 222,
	132, 54, 185, 185, 64, 134, 70, 57,
	185, 185, 222, 222, 222, 222, 222, 0,
	223, 135, 54, 186, 186, 64, 54, 70,
	57, 186, 186, 223, 223, 223, 223, 223,
	0, 224, 138, 50, 225, 129, 50, 140,
	57, 186, 186, 224, 224, 224, 224, 224,
	0, 223, 135, 54, 186, 186, 64, 54,
	70, 57, 186, 186, 223, 223, 223, 223,
	223, 0, 226, 148, 150, 54, 187, 187,
	54, 187, 57, 187, 187, 226, 226, 226,
	226, 226, 0, 227, 151, 140, 50, 228,
	50, 50, 57, 187, 187, 227, 227, 227,
	227, 227, 0, 226, 148, 150, 54, 187,
	187, 54, 187, 57, 187, 187, 226, 226,
	226, 226, 226, 0, 220, 130, 50, 221,
	129, 137, 140, 57, 230, 185, 185, 220,
	230, 220, 220, 220, 220, 0, 220, 130,
	50, 221, 129, 137, 140, 57, 231, 185,
	185, 220, 231, 220, 220, 220, 220, 0,
	220, 130, 50, 221, 129, 137, 140, 57,
	232, 185, 185, 220, 232, 220, 220, 220,
	220, 0, 220, 130, 50, 221, 129, 137,
	140, 57, 233, 185, 185, 220, 233, 220,
	220, 220, 220, 0, 220, 130, 50, 221,
	129, 137, 140, 57, 234, 185, 185, 220,
	234, 220, 220, 220, 220, 0, 220, 130,
	50, 221, 129, 153, 140, 57, 185, 185,
	220, 220, 220, 220, 220, 0, 235, 236,
	235, 50, 224, 225, 129, 50, 140, 57,
	186, 186, 196, 235, 224, 235, 235, 235,
	0, 196, 196, 196, 64, 70, 196, 196,
	196, 237, 237, 196, 196, 237, 196, 0,
	196, 196, 196, 64, 70, 196, 196, 196,
	235, 235, 196, 196, 235, 196, 0, 220,
	130, 50, 221, 129, 137, 140, 57, 239,
	185, 185, 220, 239, 220, 220, 220, 220,
	0, 220, 130, 50, 221, 129, 137, 140,
	57, 240, 185, 185, 220, 240, 220, 220,
	220, 220, 0, 220, 130, 50, 221, 129,
	137, 140, 57, 241, 185, 185, 220, 241,
	220, 220, 220, 220, 0, 220, 130, 50,
	221, 129, 137, 140, 57, 242, 185, 185,
	220, 242, 220, 220, 220, 220, 0, 220,
	130, 50, 221, 129, 137, 140, 57, 243,
	185, 185, 220, 243, 220, 220, 220, 220,
	0, 220, 130, 50, 221, 129, 137, 140,
	57, 244, 185, 185, 220, 244, 220, 220,
	220, 220, 0, 220, 130, 50, 221, 129,
	137, 140, 57, 245, 185, 185, 220, 245,
	220, 220, 220, 220, 0, 220, 130, 50,
	221, 129, 137, 140, 57, 234, 185, 185,
	220, 234, 220, 220, 220, 220, 0, 220,
	130, 50, 221, 129, 137, 140, 57, 247,
	185, 185, 220, 247, 220, 220, 220, 220,
	0, 220, 130, 50, 221, 129, 137, 140,
	57, 248, 185, 185, 220, 248, 220, 220,
	220, 220, 0, 220, 130, 50, 221, 129,
	137, 140, 57, 234, 185, 185, 220, 234,
	220, 220, 220, 220, 0, 50, 51, 124,
	126, 128, 129, 50, 140, 57, 50, 50,
	50, 250, 125, 125, 0, 50, 51, 124,
	126, 128, 129, 50, 140, 57, 50, 50,
	50, 251, 125, 125, 0, 50, 51, 124,
	126, 128, 129, 50, 140, 57, 50, 50,
	50, 125, 125, 125, 0,
}

var _uri_trans_actions []byte = []byte{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 5, 0,
	0, 0, 0, 0, 0, 0, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 0, 0, 0, 9, 0, 9,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 1, 1,
	1, 19, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 13, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 1, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 9, 0, 0, 0, 9, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	9, 0, 0, 0, 9, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	9, 0, 0, 0, 9, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 9, 0,
	0, 0, 9, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 9, 0, 0,
	0, 9, 0, 0, 0, 0, 0, 0,
	0, 1, 1, 1, 19, 1, 1, 1,
	1, 1, 1, 0, 0, 0, 0, 9,
	0, 0, 0, 9, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	13, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 9, 0, 0, 0,
	9, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 9, 0, 0,
	9, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 9, 0, 0,
	9, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 1, 1, 1, 0, 0, 1,
	0, 19, 0, 0, 1, 1, 1, 1,
	1, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 13, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 13, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 9, 0, 0, 0, 9, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 9, 0, 0, 0, 9, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 9, 0, 0, 0, 9, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 9,
	0, 0, 0, 9, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 9,
	0, 0, 0, 9, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 9, 0, 0,
	0, 9, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 9, 0, 0,
	0, 9, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 9, 0, 0,
	0, 9, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 9, 0, 0,
	0, 9, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 9, 0, 0,
	0, 9, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 9, 0, 0,
	0, 9, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 9, 0, 0,
	0, 9, 0, 0, 0, 0, 0, 0,
	0, 0, 3, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 11, 25, 34, 0, 0, 0,
	0, 11, 25, 34, 0, 0, 0, 0,
	28, 38, 0, 0, 28, 38, 0, 0,
	28, 38, 0, 0, 28, 38, 0, 0,
	28, 38, 0, 0, 0, 0, 0, 31,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 31, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 31, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 31, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 31, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 31, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	31, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 31,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 31, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 31, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 31, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	31, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 31,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 31, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 31, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 31, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 31, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 31, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 31, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	31, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 31,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 31, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 31, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 11, 25, 34, 0, 0, 0,
	0, 0, 0, 11, 25, 34, 0, 0,
	0, 0, 0, 0, 11, 25, 34, 0,
	0, 0, 0, 11, 25, 34, 0, 0,
	0, 0, 0, 22, 25, 0, 34, 9,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 22, 25, 0, 34, 9, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 28,
	0, 38, 13, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 28, 0, 38, 13,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 28, 0, 38, 13, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 28, 0,
	38, 13, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 28, 0, 38, 13, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 9, 0, 0, 31, 9, 0, 0,
	0, 0, 0, 0, 0, 0, 1, 1,
	1, 0, 0, 0, 1, 31, 19, 0,
	0, 1, 1, 1, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 31, 13,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 31,
	13, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 9, 0, 0, 31,
	9, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 1, 1, 0, 0, 0, 1,
	31, 19, 0, 0, 1, 1, 1, 1,
	1, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 13, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 9,
	0, 0, 9, 0, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 0,
	0, 1, 0, 19, 0, 0, 1, 1,
	1, 1, 1, 0, 0, 0, 0, 9,
	0, 0, 31, 9, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 9, 0, 0, 31, 9, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 9, 0, 0, 31, 9,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 9, 0, 0,
	31, 9, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 9,
	0, 0, 31, 9, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 9, 0, 0, 31, 9, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 9, 0, 0, 31, 9,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 31, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 31, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 9, 0, 0, 31, 9, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 9, 0, 0, 31,
	9, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 9, 0,
	0, 31, 9, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	9, 0, 0, 31, 9, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 9, 0, 0, 31, 9, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 9, 0, 0, 31,
	9, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 9, 0,
	0, 31, 9, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	9, 0, 0, 31, 9, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 9, 0, 0, 31, 9, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 9, 0, 0, 31,
	9, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 9, 0,
	0, 31, 9, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 22, 25, 0, 34, 9, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 22, 25, 0, 34, 9, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 22, 25, 0, 34, 9, 0, 0,
	0, 0, 0, 0, 0,
}

var _uri_to_state_actions []byte = []byte{
	0, 0, 0, 7, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0,
}

var _uri_eof_actions []byte = []byte{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 34, 34, 38, 38, 38, 38,
	38, 15, 15, 17, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 34, 34, 34, 34, 34, 34, 38,
	38, 38, 38, 38, 15, 15, 15, 15,
	15, 15, 17, 17, 17, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 34, 34, 34,
}

const uri_start int = 1
const uri_first_final int = 166
const uri_error int = 0

const uri_en_uri int = 1

//line parser_uri.rl:8

func URIParse(data []byte) *URI {
	uri := &URI{}
	uri.buf.init(data)
	cs := 0 // current state. entery point = 0
	l := uri.buf.plen()
	var p, // data pointer
		m, // marker
		pe, eof ptr = 0, 0, l, l

//line parser_uri.rl:47

//line parser_uri.go:1167
	{
		cs = uri_start
	}

//line parser_uri.rl:49

//line parser_uri.go:1174
	{
		var _klen int
		var _trans int
		var _acts int
		var _nacts uint
		var _keys int
		if p == pe {
			goto _test_eof
		}
		if cs == 0 {
			goto _out
		}
	_resume:
		_keys = int(_uri_key_offsets[cs])
		_trans = int(_uri_index_offsets[cs])

		_klen = int(_uri_single_lengths[cs])
		if _klen > 0 {
			_lower := int(_keys)
			var _mid int
			_upper := int(_keys + _klen - 1)
			for {
				if _upper < _lower {
					break
				}

				_mid = _lower + ((_upper - _lower) >> 1)
				switch {
				case data[p] < _uri_trans_keys[_mid]:
					_upper = _mid - 1
				case data[p] > _uri_trans_keys[_mid]:
					_lower = _mid + 1
				default:
					_trans += int(_mid - int(_keys))
					goto _match
				}
			}
			_keys += _klen
			_trans += _klen
		}

		_klen = int(_uri_range_lengths[cs])
		if _klen > 0 {
			_lower := int(_keys)
			var _mid int
			_upper := int(_keys + (_klen << 1) - 2)
			for {
				if _upper < _lower {
					break
				}

				_mid = _lower + (((_upper - _lower) >> 1) & ^1)
				switch {
				case data[p] < _uri_trans_keys[_mid]:
					_upper = _mid - 2
				case data[p] > _uri_trans_keys[_mid+1]:
					_lower = _mid + 2
				default:
					_trans += int((_mid - int(_keys)) >> 1)
					goto _match
				}
			}
			_trans += _klen
		}

	_match:
		cs = int(_uri_trans_targs[_trans])

		if _uri_trans_actions[_trans] == 0 {
			goto _again
		}

		_acts = int(_uri_trans_actions[_trans])
		_nacts = uint(_uri_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _uri_actions[_acts-1] {
			case 0:
//line parser_uri.rl:18
				m = p
			case 1:
//line parser_uri.rl:19
				uri.scheme = pl{0, p}
				uri.id = URIsips
			case 2:
//line parser_uri.rl:20
				uri.scheme = pl{0, p}
				uri.id = URIsip
			case 4:
//line parser_uri.rl:22

				from := uri.scheme.l + 1
				if uri.id == URIabs {
					from = m
				}
				uri.user = pl{from, p}

			case 5:
//line parser_uri.rl:27
				uri.host = pl{m, p}
			case 6:
//line parser_uri.rl:28
				uri.password = pl{m, p}
			case 7:
//line parser_uri.rl:29
				uri.port = pl{m, p}
			case 8:
//line parser_uri.rl:30
				uri.params = pl{m, p}
//line parser_uri.go:1280
			}
		}

	_again:
		_acts = int(_uri_to_state_actions[cs])
		_nacts = uint(_uri_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _uri_actions[_acts-1] {
			case 3:
//line parser_uri.rl:21
				uri.scheme = pl{0, p}
				uri.id = URIabs
//line parser_uri.go:1293
			}
		}

		if cs == 0 {
			goto _out
		}
		p++
		if p != pe {
			goto _resume
		}
	_test_eof:
		{
		}
		if p == eof {
			__acts := _uri_eof_actions[cs]
			__nacts := uint(_uri_actions[__acts])
			__acts++
			for ; __nacts > 0; __nacts-- {
				__acts++
				switch _uri_actions[__acts-1] {
				case 0:
//line parser_uri.rl:18
					m = p
				case 5:
//line parser_uri.rl:27
					uri.host = pl{m, p}
				case 7:
//line parser_uri.rl:29
					uri.port = pl{m, p}
				case 8:
//line parser_uri.rl:30
					uri.params = pl{m, p}
				case 9:
//line parser_uri.rl:31
					uri.headers = pl{m, p}
//line parser_uri.go:1326
				}
			}
		}

	_out:
		{
		}
	}

//line parser_uri.rl:50
	if cs >= uri_first_final {
		return uri
	}
	return nil
}
