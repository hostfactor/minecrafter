mocks:
	for i in `find . -name mock_*`; do rm -f $i; done
	mockery --all --dir . --case snake
