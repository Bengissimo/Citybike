#!/bin/bash

links=(
"http://localhost:8000" #200
"http://localhost:8000/journeys" #200
"http://localhost:8000/journeys?p=2" #200
"http://localhost:8000/journeys?p=100000" #307
"http://localhost:8000/journeys?p=abc" #400
"http://localhost:8000/stations" #200
"http://localhost:8000/stations?p=2" #200
"http://localhost:8000/stations?p=100" #307
"http://localhost:8000/stations?p=abc" #400
"http://localhost:8000/station/1" #200
"http://localhost:8000/station/" #400
"http://localhost:8000/station/1000" #404
"http://localhost:8000/station/abc" #400
)

expected_codes=(200 200 200 307 400 200 200 307 400 200 400 404 400)

length=${#links[@]}

for ((i=0; i<$length; i++)); do
    link=${links[i]}
    expected_code=${expected_codes[i]}

    # Send a GET request to the link using curl
    status_code=$(curl -w "%{http_code}\n" -o /dev/null -s $link)

    if [ $status_code -eq $expected_code ]; then
        echo -e "\033[32mOK\033[0m $link : $expected_code and $status_code"
    else
        echo -e "\033[31mERROR\033[0m $link : $expected_code expected, but $status_code received."
    fi
done

echo "-------"

firstpages=("http://localhost:8000/stations?p=0" "http://localhost:8000/journeys?p=0")

for link in "${firstpages[@]}"
do
	output=$(curl -s $link)
	substring=$(basename $link | cut -d "?" -f 1)

# Check if the string "Previous" appears in the first page
	if echo "$output" | grep -q "Previous"; then
    	echo -e "\033[31mERROR\033[0m In $substring 'Previous' found in the first page"
	else
		if echo "$output" | grep -q "Next"; then
			echo -e "\033[32mOK\033[0m 'Previous' not found and 'Next' found in the first page of $substring list"
		else
			echo -e "\033[32mERROR\033[0m No Previous No Next link in the first page of $substring list: "
		fi
	fi
done

echo "-------"

lastpages=("http://localhost:8000/stations?p=22" "http://localhost:8000/journeys?p=58159")

for link in "${lastpages[@]}"
do
	output=$(curl -s $link)
	substring=$(basename $link | cut -d "?" -f 1)

# Check if the string "Next" appears in the last page
	if echo "$output" | grep -q "Next"; then
    	echo -e "\033[31mERROR\033[0m In $substring 'Next' found in the last page"
	else
		if echo "$output" | grep -q "Previous"; then
			echo -e "\033[32mOK\033[0m 'Next' not found and 'Previous' found in the last page of $substring list"
		else
			echo -e "\033[31mERROR\033[0m No Previous No Next link in the last page of $substring list: "
		fi
	fi
done

echo "-------"

midpages=("http://localhost:8000/stations?p=2" "http://localhost:8000/journeys?p=2")

for link in "${midpages[@]}"
do
	output=$(curl -s $link)
	substring=$(basename $link | cut -d "?" -f 1)

# Check if the string "Next" appears in the last page
	if echo "$output" | grep -q "Next\|Previous"; then
    	echo -e "\033[32mOK\033[0m 'Next' and 'Previous' found in the second page of $substring "
	else
		echo -e "\033[31mERROR\033[0m 'Next' not found and 'Previous' found in the second page of $substring list"
	fi
done
