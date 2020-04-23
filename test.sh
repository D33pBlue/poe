# @Author: Bizzaro Francesco <d33pblue>
# @Date:   2020-Apr-23
# @Project: Proof of Evolution
# @Filename: test.sh
# @Last modified by:   d33pblue
# @Last modified time: 2020-Apr-23
# @Copyright: 2020
for D in `find -maxdepth 1 -type d ! -regex '\.\(/\..*\)?'`
do
    cd $D
    nn=`find -maxdepth 1 -regex '.*_test.go'`
    xyz="$nn"; if [ -z "$xyz" ] ; then R=""; else go test; fi
    cd ..
done
