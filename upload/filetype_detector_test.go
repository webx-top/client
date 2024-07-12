package upload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var svgBytes = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg width="464" height="153" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
<!--WillInclude_xlog-->
<path d="M 124.23 49.94 C 130.58 56.45 130.58 67.02 124.23 73.53 C 117.88 80.05 107.59 80.05 101.24 73.53 C 94.89 67.02 94.89 56.45 101.24 49.94 C 107.59 43.42 117.88 43.42 124.23 49.94" stroke-width="1.00" stroke="white" stroke-linecap="round" stroke-linejoin="round" fill="#fc7708"  />

<path d="M 182.16 13.91 C 186.47 18.04 186.47 24.74 182.16 28.88 C 177.84 33.01 170.85 33.01 166.53 28.88 C 162.22 24.74 162.22 18.04 166.53 13.91 C 170.85 9.78 177.84 9.78 182.16 13.91" stroke-width="1.00" stroke="white" stroke-linecap="round" stroke-linejoin="round" fill="#f47408"  />

<path d="M 150.16 24.12 C 155.00 28.97 155.00 36.82 150.16 41.66 C 145.32 46.50 137.47 46.50 132.63 41.66 C 127.79 36.82 127.79 28.97 132.63 24.12 C 137.47 19.28 145.32 19.28 150.16 24.12" stroke-width="1.00" stroke="white" stroke-linecap="round" stroke-linejoin="round" fill="#f47408"  />

<path d="M 213.27 23.36 C 217.59 27.49 217.59 34.19 213.27 38.32 C 208.96 42.46 201.96 42.46 197.65 38.32 C 193.33 34.19 193.33 27.49 197.65 23.36 C 201.96 19.22 208.96 19.22 213.27 23.36" stroke-width="1.00" stroke="white" stroke-linecap="round" stroke-linejoin="round" fill="#f47408"  />

<path d="M 227.37 54.25 C 231.69 58.39 231.69 65.09 227.37 69.22 C 223.06 73.35 216.06 73.35 211.75 69.22 C 207.43 65.09 207.43 58.39 211.75 54.25 C 216.06 50.12 223.06 50.12 227.37 54.25" stroke-width="1.00" stroke="white" stroke-linecap="round" stroke-linejoin="round" fill="#f47408"  />

<path d="M 382.61 128.34 L 389.83 128.34 L 389.83 135.11 L 382.61 135.11 Z M 402.05 110.79 L 402.05 117.49 L 405.77 117.49 L 405.77 122.42 L 402.05 122.42 L 402.05 128.68 C 402.05 129.43 402.12 129.93 402.27 130.17 C 402.49 130.54 402.88 130.73 403.43 130.73 C 403.93 130.73 404.62 130.59 405.52 130.30 L 406.02 134.96 C 404.35 135.33 402.79 135.51 401.34 135.51 C 399.66 135.51 398.42 135.29 397.62 134.86 C 396.83 134.43 396.24 133.78 395.86 132.90 C 395.47 132.02 395.28 130.59 395.28 128.63 L 395.28 122.42 L 392.80 122.42 L 392.80 117.49 L 395.28 117.49 L 395.28 114.26 Z M 408.16 126.35 C 408.16 123.67 409.06 121.45 410.88 119.71 C 412.69 117.97 415.14 117.10 418.23 117.10 C 421.75 117.10 424.42 118.12 426.22 120.17 C 427.67 121.81 428.39 123.84 428.39 126.25 C 428.39 128.96 427.50 131.18 425.70 132.91 C 423.90 134.64 421.42 135.51 418.24 135.51 C 415.41 135.51 413.12 134.79 411.37 133.35 C 409.23 131.57 408.16 129.24 408.16 126.35 Z M 414.92 126.34 C 414.92 127.91 415.24 129.07 415.88 129.82 C 416.51 130.57 417.31 130.95 418.28 130.95 C 419.25 130.95 420.05 130.58 420.67 129.84 C 421.30 129.10 421.61 127.91 421.61 126.27 C 421.61 124.74 421.29 123.61 420.66 122.86 C 420.03 122.11 419.25 121.74 418.33 121.74 C 417.34 121.74 416.53 122.12 415.89 122.88 C 415.25 123.64 414.92 124.79 414.92 126.34 Z M 431.68 141.81 L 431.68 117.49 L 438.00 117.49 L 438.00 120.10 C 438.87 119.00 439.67 118.26 440.40 117.88 C 441.39 117.36 442.48 117.10 443.67 117.10 C 446.03 117.10 447.85 118.00 449.14 119.80 C 450.43 121.60 451.07 123.83 451.07 126.49 C 451.07 129.42 450.37 131.65 448.96 133.20 C 447.56 134.74 445.78 135.51 443.64 135.51 C 442.60 135.51 441.65 135.33 440.79 134.98 C 439.94 134.63 439.17 134.10 438.50 133.40 L 438.50 141.81 Z M 438.45 126.35 C 438.45 127.75 438.74 128.78 439.33 129.46 C 439.91 130.13 440.65 130.47 441.55 130.47 C 442.33 130.47 442.99 130.14 443.52 129.50 C 444.05 128.85 444.32 127.75 444.32 126.20 C 444.32 124.78 444.04 123.73 443.49 123.06 C 442.94 122.39 442.26 122.06 441.47 122.06 C 440.60 122.06 439.88 122.39 439.31 123.07 C 438.73 123.74 438.45 124.84 438.45 126.35 Z M 452.26 135.11" stroke-width="1.00" stroke="white" stroke-linecap="round" stroke-linejoin="round" fill="#fc7808"  />

<path d="M 12.59 57.15 L 36.63 57.15 L 36.63 84.38 L 62.90 84.38 L 62.90 57.15 L 87.04 57.15 L 87.04 134.95 L 62.90 134.95 L 62.90 103.48 L 36.63 103.48 L 36.63 134.95 L 12.59 134.95 L 12.59 57.15 Z M 101.79 84.38 L 123.39 84.38 L 123.39 134.95 L 101.79 134.95 L 101.79 84.38 Z M 138.14 57.15 L 197.58 57.15 L 197.58 73.87 L 162.29 73.87 L 162.29 87.46 L 192.43 87.46 L 192.43 103.16 L 162.29 103.16 L 162.29 134.95 L 138.14 134.95 L 138.14 57.15 Z M 209.30 78.59 L 230.90 78.59 L 230.90 134.95 L 209.30 134.95 L 209.30 78.59 Z M 244.85 57.15 L 266.50 57.15 L 266.50 134.95 L 244.85 134.95 L 244.85 57.15 Z M 279.87 78.59 L 300.04 78.59 L 300.04 86.82 C 302.94 83.39 305.87 80.95 308.82 79.50 C 311.77 78.04 315.34 77.32 319.51 77.32 C 324.01 77.32 327.56 78.12 330.18 79.71 C 332.80 81.30 334.94 83.67 336.60 86.82 C 340.00 83.14 343.09 80.64 345.89 79.31 C 348.68 77.98 352.13 77.32 356.24 77.32 C 362.29 77.32 367.01 79.12 370.41 82.71 C 373.80 86.30 375.50 91.91 375.50 99.56 L 375.50 134.95 L 353.85 134.95 L 353.85 102.85 C 353.85 100.30 353.35 98.41 352.36 97.17 C 350.91 95.22 349.11 94.25 346.95 94.25 C 344.40 94.25 342.35 95.17 340.80 97.01 C 339.24 98.85 338.46 101.80 338.46 105.87 L 338.46 134.95 L 316.81 134.95 L 316.81 103.91 C 316.81 101.43 316.67 99.75 316.38 98.87 C 315.92 97.45 315.13 96.31 313.99 95.44 C 312.86 94.58 311.54 94.14 310.01 94.14 C 307.54 94.14 305.50 95.08 303.91 96.96 C 302.32 98.83 301.52 101.91 301.52 106.19 L 301.52 134.95 L 279.87 134.95 L 279.87 78.59 Z M 279.87 78.59" stroke-width="1.00" stroke="white" stroke-linecap="round" stroke-linejoin="round" fill="#fc7808"  />


</svg>
<!--OldSVGSize:5491 -->
<!--Data_xlog_bIncluded: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-->
`)

func TestIsSVG(t *testing.T) {
	r := IsSVGImage(svgBytes)
	assert.True(t, r)
}
