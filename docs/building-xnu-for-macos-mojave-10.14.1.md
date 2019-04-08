# Building XNU for macOS Mojave 10.14.1

> https://kernelshaman.blogspot.com/2018/12/building-xnu-for-macos-mojave-10141.html

The macOS Mojave 10.14.1 kernel (XNU) source has been released here: [source](https://opensource.apple.com/source/xnu/xnu-4903.221.2/)[](https://www.blogger.com/), [tarball](https://opensource.apple.com/tarballs/xnu/xnu-4903.221.2.tar.gz).  
  
Building XNU requires some patience, and some open source dependencies which are not pre-installed. This post walks through all the steps necessary to build the open source version of XNU on supported Apple hardware.  
  
~~**WARNING: Unfortunately, the open source Mojave 10.14.1 kernel does not include all symbols necessary to successfully create a prelinkedkernel image. This means that open source XNU will not boot, even in a VM.**~~
  
**UPDATE!**  
**The [Makefile](Makefile.xnudeps) (https://jeremya.com/sw/Makefile.xnudeps) has been updated to provide a set of fixups that allow the open source XNU kernel to be built and run. In order to do this, you must also build a new System.kext. Doing this may render your previous kernel un-bootable - proceed with caution!**
  

## TL;DR

There is a Makefile which automates the downloading and building of all prerequisites. You can find the Makefile [here](Makefile.xnudeps) (https://jeremya.com/sw/Makefile.xnudeps), and invoke it like:  

```sh
make -f Makefile.xnudeps xnudeps
```

NEW: this Makefile will now automatically detect the correct versions of source code to download based on the version of macOS you specify. By default, the version is 10.14.1, however you can select a different version like:  

```sh
make -f Makefile.xnudeps macos\_version=10.13.1 xnudeps
```

You can also see other features of the Makefile using the help target. Note that full 10.13.x compilation support will be coming soon.  
  
**UPDATE:**  
**the default target, `xnudeps`, will perform all the necessary source fixups. However, if you want to re-download the source and perform the fixups, you can use the `xnu-fixups` target.**

## Manual XNU Building

All of the source for both XNU and required dependencies is available from [opensource.apple.com](https://opensource.apple.com/). Here are the manual steps necessary to build XNU:  

### 1. Download and Install Xcdoe

- Make sure you have Xcode 10 (or 10.1) installed. You can install it via the App Store, or by manual download here:
    - [https://developer.apple.com/download/more/](https://developer.apple.com/download/more/)
- NOTE: for older versions of macOS, you may need older versions of Xcode

### 2. Download the source

```sh
export TARBALLS=https://opensource.apple.com/tarballs
curl -O ${TARBALLS}/dtrace/dtrace-284.200.15.tar.gz
curl -O ${TARBALLS}/AvailabilityVersions/AvailabilityVersions-33.200.4.tar.gz
curl -O ${TARBALLS}/libplatform/libplatform-177.200.16.tar.gz
curl -O ${TARBALLS}/libdispatch/libdispatch-1008.220.2.tar.gz
curl -O ${TARBALLS}/xnu/xnu-4903.221.2.tar.gz
```

### 3. Build CTF tools from dtrace

```sh
tar zxf dtrace-284.200.15.tar.gz
cd dtrace-284.200.15
mkdir obj sym dst
xcodebuild install -sdk macosx -target ctfconvert -target ctfdump -target ctfmerge ARCHS=x86_64 SRCROOT=$PWD OBJROOT=$PWD/obj SYMROOT=$PWD/sym DSTROOT=$PWD/dst HEADER_SEARCH_PATHS="$PWD/compat/opensolaris/** $PWD/lib/**"
sudo ditto $PWD/dst/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain
cd ..
```

### 4. Install AvailabilityVersions

```sh
tar zxf AvailabilityVersions-33.200.4.tar.gz
cd AvailabilityVersions-33.200.4
mkdir dst
make install SRCROOT=$PWD DSTROOT=$PWD/dst
sudo ditto $PWD/dst/usr/local/libexec $(xcrun -sdk macosx -show-sdk-path)/usr/local/libexec
cd ..
```

### 5. Install libplatform headers

```sh
tar zxf libplatform-177.200.16.tar.gz
cd libplatform-177.200.16
sudo mkdir -p $(xcrun -sdk macosx -show-sdk-path)/usr/local/include/os/internal
sudo ditto $PWD/private/os/internal $(xcrun -sdk macosx -show-sdk-path)/usr/local/include/os/internal
cd ..
```

### 6. Install XNU headers

```sh
tar zxf xnu-4903.221.2.tar.gz
cd xnu-4903.221.2
make SDKROOT=macosx ARCH_CONFIGS=X86_64 installhdrs
sudo ditto $PWD/BUILD/dst $(xcrun -sdk macosx -show-sdk-path)
cd ..
```

### 7. Build firehose from libdispatch

```sh
tar zxf libdispatch-1008.220.2.tar.gz
cd libdispatch-1008.220.2
mkdir obj sym dst
awk '/include "<DEVELOPER/ {next;} /SDKROOT =/ {print "SDKROOT = macosx"; next;} {print $0}' xcodeconfig/libdispatch.xcconfig > .__tmp__ && mv -f .__tmp__ xcodeconfig/libdispatch.xcconfig
awk '/#include / { next; } { print $0 }' \xcodeconfig/libfirehose_kernel.xcconfig > .__tmp__ && mv -f .__tmp__ xcodeconfig/libfirehose_kernel.xcconfig
xcodebuild install -sdk macosx -target libfirehose_kernel SRCROOT=$PWD OBJROOT=$PWD/obj SYMROOT=$PWD/sym DSTROOT=$PWD/dst
sudo ditto $PWD/dst/usr/local $(xcrun -sdk macosx -show-sdk-path)/usr/local
cd ..
```

### 8. Build XNU (checkout the README.md for more options!)

```sh
cd xnu-4903.221.2
make SDKROOT=macosx ARCH_CONFIGS=X86_64 KERNEL_CONFIGS=RELEASE
```

Check out the [README.md](https://github.com/apple/darwin-xnu/blob/master/README.md) file at the top of the XNU source tree for more options to the build system. Some common and useful options include:

```console
KERNEL_CONFIGS=DEVELOPMENT BUILD_LTO=n LOGCOLORS=y
```

### Install and Run XNU

NOTE: You may need to [disable System Integrity Protection](https://developer.apple.com/library/archive/documentation/Security/Conceptual/System_Integrity_Protection_Guide/ConfiguringSystemIntegrityProtection/ConfiguringSystemIntegrityProtection.html) in order to install and run a custom kernel.

~~**WARNING: Unfortunately, the open source Mojave 10.14.1 kernel does not include all symbols necessary to successfully create a prelinkedkernel image. This means that open source XNU will not boot, even in a VM.**~~
  
**WARNING: You will need to perform some source fixups on the xnu open source drop. If you have manually followed these steps, you can perform the fixups by doing this:**

- `curl -O https://jeremya.com/sw/Makefile.xnudeps`
- `make -f Makefile.xnudeps xnu-fixups`
- {rebuild XNU}

  
After the final build step, you should have a new kernel built in `$PWD/BUILD/obj/kernel`. In order to run this kernel, you will need to install it, and rebuild the prelinkedkernel image. Installing a kernel could potentially render your system un-bootable, so trying this out in a VM first is recommended.  
  
In order to successfully link the macOS Mojave open source kernel, you will need to build and install the System.kext. Fortunately, this is straightforward:  

```sh
# make a backup copy of the existing System.kext!
sudo ditto /System/Library/Extensions/System.kext ~/System.kext.backup
cd xnu-4903.221.2
make SDKROOT=macosx KERNEL_CONFIGS=RELEASE DSTROOT=$PWD/BUILD.syskext \install\_config
sudo chown -R root:wheel BUILD.syskext
sudo ditto BUILD.syskext/ /
cd ..
```
  
To install and run your kernel:  

```sh
cd xnu-4903.221.2
sudo ditto $PWD/BUILD/obj/kernel /System/Library/Kernels/kernel
sudo kextcache -v -invalidate /
# / locked; waiting for lock.
# Lock acquired; proceeding
# ...

sudo reboot
# ...

uname -a
```

If you build a different variant of XNU, you may need to ditto a different kernel name, e.g., `kernel.development` instead of just `kernel`.

Note that you can select different prelinkedkernel variants from which to boot using the kcsuffix boot-arg. For example, if you built a development kernel (using `KERNEL_CONFIGS=DEVELOPMENT` in the make invocation), you would install and run it like so:

```sh
sudo ditto $PWD/BUILD/obj/kernel.development /System/Library/Kernels/kernel.development
sudo kextcache -v -invalidate /
sudo nvram boot-args="kcsuffix=development"
sudo reboot
```

If you have existing boot-args, you can, of course, preserve them in the nvram boot-args variable.
