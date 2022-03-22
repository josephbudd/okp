# How I created this assets folder.

1.  I started here at github.com.
    I created the remote assets branch in my okp project at github.com/josephbudd/okp.
    The assets branch ended up being a copy of the master with all of its files.


2.  I cloned the branch locally.
    I used the git branch command to confirm that this is the assets branch.

nil@NIL:~/go/src/github.com/josephbudd/okp_assets$ git clone --branch assets https://BLAH_BLAHBLAHBLAHBLAHBLAH@github.com/josephbudd/okp_assets.git
Cloning into 'okp'...
 ...
nil@NIL:~/go/src/github.com/josephbudd/okp_assets$ cd okp
nil@NIL:~/go/src/github.com/josephbudd/okp_assets/okp$ git branch
* assets


3.a I deleted files and folders from the local branch.
    I committed and pushed the deletes up to the remote branch.

nil@NIL:~/go/src/github.com/josephbudd/okp_assets/okp$ git add .
nil@NIL:~/go/src/github.com/josephbudd/okp_assets/okp$ git commit -m "removed files"
 ...
nil@NIL:~/go/src/github.com/josephbudd/okp_assets/okp$ git push https://BLAH_BLAHBLAHBLAHBLAHBLAH@github.com/josephbudd/okp_assets.git
 ...

3.b I looked at github.com/josephbudd/okp assets branch and there are no more files. 


4.a I copied my 2 image files into the local branch and committed and pushed the additions up to the remote branch.

nil@NIL:~/go/src/github.com/josephbudd/okp_assets/okp$ git add .
nil@NIL:~/go/src/github.com/josephbudd/okp_assets/okp$ git commit -m "added image files"
 ...
nil@NIL:~/go/src/github.com/josephbudd/okp_assets/okp$ git push https://BLAH_BLAHBLAHBLAHBLAHBLAH@github.com/josephbudd/okp_assets.git
 ...


4.b I looked at github.com/josephbudd/okp assets branch and there are the 2 new image files.
    I looked at github.com/josephbudd/okp master branch and it is untouched.


