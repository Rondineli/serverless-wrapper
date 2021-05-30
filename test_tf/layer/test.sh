				mkdir -p "./python/lib/python3.8/site-packages/" || echo "ok"
				cp *.py "./python/lib/python3.8/site-packages/" || echo "ok"

				if [ "pip" == "pip" ]; then
					python3 -m pip install -r ./requirements.txt -t "./python/lib/python3.8/site-packages"
				else
					## If not pip only should be something like pip3.8, pip3.7 etc..
					pip3 install -r ./requirements.txt -t "./python/lib/python3.8/site-packages"
				fi
