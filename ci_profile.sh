if [[ $TRAVIS_BRANCH = 'master' ]]; then
    export DHOST='47.100.0.27'
else
    export DHOST='101.132.143.0'
fi

export DPORT='2022'

echo $TRAVIS_BRANCH
echo $DPORT
echo $DHOST