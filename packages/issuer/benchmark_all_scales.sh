#!/bin/bash

echo "=========================================="
echo "MULTI-SCALE BENCHMARK TEST"
echo "=========================================="
echo ""

scales=(100 500 1000 2500 5000)

for students in "${scales[@]}"; do
    echo "=========================================="
    echo "TESTING: $students STUDENTS"
    echo "=========================================="
    
    # Reset
    ./reset.sh > /dev/null 2>&1
    
    # Generate data
    cd cmd
    echo "Generating data for $students students..."
    time_output=$(time -p go run . generate-data --students=$students 2>&1)
    gen_time=$(echo "$time_output" | grep "^real" | awk '{print $2}')
    
    # Count courses
    total_courses=$(echo "$time_output" | grep "Total Course Completions" | awk '{print $NF}')
    
    echo "  ✅ Generated: $students students, $total_courses courses in ${gen_time}s"
    
    # Build one term tree
    echo "Building Verkle tree for Semester_1_2023..."
    go run . convert-data Semester_1_2023 > /dev/null 2>&1
    
    tree_start=$(date +%s.%N)
    tree_output=$(go run . add-term Semester_1_2023 ../data/verkle_terms/Semester_1_2023_completions.json 2>&1)
    tree_end=$(date +%s.%N)
    tree_time=$(echo "$tree_end - $tree_start" | bc)
    
    tree_courses=$(echo "$tree_output" | grep "Loaded" | awk '{print $3}')
    
    echo "  ✅ Verkle tree: $tree_courses courses in ${tree_time}s"
    
    # Generate one receipt
    echo "Generating receipt for ITITIU00001..."
    receipt_start=$(date +%s.%N)
    go run . generate-receipt ITITIU00001 ../publish_ready/receipts/test_${students}.json > /dev/null 2>&1
    receipt_end=$(date +%s.%N)
    receipt_time=$(echo "$receipt_end - $receipt_start" | bc)
    
    receipt_size=$(ls -l ../publish_ready/receipts/test_${students}.json 2>/dev/null | awk '{print $5}')
    
    echo "  ✅ Receipt: ${receipt_size} bytes in ${receipt_time}s"
    
    echo ""
    cd ..
done

echo "=========================================="
echo "BENCHMARK COMPLETE"
echo "=========================================="
