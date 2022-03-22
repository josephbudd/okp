# okp "One Key Piano"

## This is a work in progress

## About okp

okp is my morse code trainer that I developed to teach myself morse code.

Back in the 90's, I used a cassette tape and a straight key to learn morse code. A few years ago I wrote my own Go code to translate my keying using a browser and a USB key. That's when I realized that I was not keying morse code, I was keying a bunch of noise.

Morse code is not noise. Morse code is music.

## At this state of development, the application has 2 parts, "Courses" and "Training"

The application works fine. I'm testing for little things that I may have missed.

### 1. Courses

In **Courses** one can

* View a description of the current course and select a different current course.
* Create a new course.
* Edit an existing course.
* Remove an existing course.

#### A course has

* A name.
* A description.
* One of the various lesson plans. A lesson plan is series of lessons. A lesson presents the character, word or sentence that one must learn to copy and key.
* One of the various speeds for keying and copying.

### 2. Training

In **Training** one learns the current course one lesson at a time. In a lesson one pratices as long as desired and then tests.

#### Copying

##### Copy Practice

* The morse code is keyed by the app.
* Correct copies don't count toward anything. Incorrect copies don't count against anything.

##### Copy Test

* The morse code is keyed by the app.
* Correct copy attempts accumulate until the required amount of correct copy attempts is reached.
* Incorrect copy attempts have no effect on the accumulated correct copy attempts.

#### Keying

##### Key Practice

* The text to be keyed is displayed.
* The keying instructions are displayed.
* The app's metronome can be turned on to keep time.
* When the metronome is turned off the key sounds in the app. The key has a subtle beat to help keep time. Hold the key down for 1 beat for a dit. Hold the key down for 3 beats for a dah.
* Correct key attempts don't count toward anything.
* Incorrect key attempts don't count against anything.

##### Key Test

* The text to be keyed is displayed.
* The keying instructions are not displayed.
* The app's metronome can not be turned on to keep time.
* The key sounds in the app. The key has a subtle beat to help keep time. Hold the key down for 1 beat for a dit. Hold the key down for 3 beats for a dah.
* Correct key attempts accumulate until the required amount of correct key attempts is reached.
* Incorrect key attempts have no effect on the accumulated correct key attempts.

## Data stores.

The application stores it's data in easy to read text files at ~/.okp/stores/*.yaml. The **~/.okp/** folder can be deleted at any time while the app is not running.

## How to build okp

I'm building this application on ubuntu 20 using Go, CGO, using Go Modules, Go Workspaces and VSCode workspaces.

I am using the [Fyne](https://fyne.io/) GUI. The Fyne GUI is made for all devices so it's widgets work on all devices.

I wanted a simple rectangular widget with call backs for mouse events. For this reason, I used a Go Workspace because it allows me to have access to the fyne widget folder which is normally read only.

### Step 1: Create the Go Workspace at ~/workspace_okp/

```shell
nil@NIL:~$ cd
nil@NIL:~$ mkdir workspace_okp
nil@NIL:~$ cd workspace_okp/
nil@NIL:~/workspace_okp$ git clone https://github.com/josephbudd/okp
nil@NIL:~/workspace_okp$ git clone https://github.com/fyne-io/fyne.git
nil@NIL:~/workspace_okp$ go work init
nil@NIL:~/workspace_okp$ go work use ./fyne
nil@NIL:~/workspace_okp$ go work use ./okp
```

### Step 2: Install the libasound2 library which is for the okp/backend/model/goalsa package

ALSA is the Advanced Linux Sound Architecture.

libasound2 is needed for the okp/backend/model/goalsa package. libasound2 is the shared library for ALSA applications. It contains the ALSA library and its standard plugins, as well as the required configuration files.

```shell
nnil@NIL:~/workspace_okp$ sudo apt update
nil@NIL:~/workspace_okp$ sudo apt install libasound2-dev
```

### Step 3: Add my mousepad widget to fyne

The source code for my mouse pad widget needs to be copied from a hidden folder in okp to where it belongs in the widget folder in fyne.

```shell
nil@NIL:~/workspace_okp$ cp ./okp/_files/mousepad.go ./fyne/widget/
```

### Step 4: Build and run okp

**main.go** is in the ./okp/shared folder so we will build there and then put the executable in the ./okp folder.

```shell
nil@NIL:~/workspace_okp$ go build -o ./okp/okp ./okp/shared
nil@NIL:~/workspace_okp$ ./okp/okp
```

## How to uninstall okp.

You don't have to uninstall okp because there is nothing to uninstall. You only have to delete

1. the executable file **okp** from where ever you put it.
1. the **.okp/** data folder at **~/.okp/**

## How to edit okp with vscode

### 1 editor

```shell
nil@NIL:~/workspace_okp$ code ./okp
```

### 2 editors

1. An editor for the front end which also edits the shared folder.
1. An editor for the back end which also edits the shared folder.


```shell
nil@NIL:~/workspace_okp$ code ./okp/backend.code-workspace 
nil@NIL:~/workspace_okp$ code ./okp/frontend.code-workspace 
```

## How to build a USB straigt key.

![Mouse](../assets/USBKeyMouse.png)

1. Get a mouse that you are willing to tear apart.
2. Expose the left mouse button.
3. Solder wires to each of the 2 contacts on the left mouse button.

![Key](../assets/USBKey.png)

4. Attach the other end of the 2 wires to the straight key.
5. Plug it into an empty USB port on you computer.
