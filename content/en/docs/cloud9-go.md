---
title: "Using AWS Cloud9 with the AWS SDK for Go V2"
linkTitle: "Using Cloud9 with the SDK"
date: "2020-11-12"
---

You can use {{% alias service=AC9long %}} with the {{% alias sdk-go %}} to write and run your Go code using just a
browser. {{% alias service=AC9 %}} includes tools such as a code editor and terminal. Because the {{% alias service=AC9
%}} IDE is cloud based, you can work on your projects from your office, home, or anywhere using an internet-connected
machine. For general information about {{% alias service=AC9 %}}, see the {{% alias service=AC9 %}} User Guide.

Follow these instructions to set up {{% alias service=AC9 %}} with the {{% alias sdk-go %}}:

## Set up Your AWS Account to Use {{% alias service=AC9 %}}

Start to use {{% alias service=AC9 %}} by signing in to the {{% alias service=AC9 %}} console as an {{% alias
service=IAMlong %}} ({{% alias service=IAM %}}) entity (for example, an {{% alias service=IAM %}} user) in your AWS
account who has access permissions for {{% alias service=AC9 %}}.

To set up an {{% alias service=IAM %}} entity in your AWS account to access {{% alias service=AC9 %}}, and to sign in to
the {{% alias service=AC9 %}} console, see
[Team Setup for AWS Cloud9](https://docs.aws.amazon.com/cloud9/latest/user-guide/setup.html) in the {{% alias
service=AC9 %}} User Guide.

## Set up Your {{% alias service=AC9 %}} Development Environment

After you sign in to the {{% alias service=AC9 %}} console, use the console to create an {{% alias service=AC9 %}}
development environment. After you create the environment, {{% alias service=AC9 %}} opens the IDE for that environment.

See [Creating an Environment in AWS Cloud9](https://docs.aws.amazon.com/cloud9/latest/user-guide/create-environment.html)
in the {{% alias service=AC9 %}} User Guide for details.

{{% pageinfo color="info" %}} As you create your environment in the console for the first time, we recommend that you
choose the option to **Create a new instance for environment (EC2)**. This option tells {{% alias service=AC9 %}} to
create an environment, launch an {{% alias service=EC2 %}} instance, and then connect the new instance to the new
environment. This is the fastest way to begin using {{% alias service=AC9 %}}. {{% /pageinfo %}}

## Set up the AWS SDK for Go V2

After {{% alias service=AC9 %}} opens the IDE for your development environment, use the IDE to set up the {{% alias
sdk-go %}} in your environment, as follows.

1. If the terminal isn't already open in the IDE, open it. On the menu bar in the IDE, choose **Window, New Terminal**.
1. Set your GOPATH environment variable. To do this, add the following code to the end of your shell profile file (for
   example, :file:`~/.bashrc` in Amazon Linux, assuming you chose the option to **Create a new instance for
   environment (EC2)**, earlier in this topic), and then save the file.
   ```bash
   GOPATH=~/environment/go
   export GOPATH
   ```
   After you save the file, source the :file:`~/.bashrc` file to finish setting your GOPATH environment variable. To do
   this, run the following command. (This command assumes you chose the option to **Create a new instance for
   environment (EC2)**, earlier in this topic.)
   ```bash
   . ~/.bashrc
   ```
1. Run the following command to install the {{% alias sdk-go %}}.
   ```bash
   go get -u github.com/aws/aws-sdk-go/...
   ```

If the IDE can't find Go, run the following commands, one at a time in this order, to install it. (These commands assume
you chose the option to **Create a new instance for environment (EC2)**, earlier in this topic. Also, these commands
assume the latest stable version of Go at the time this topic was written; for more information,
see [Downloads](https://golang.org/dl/) on The Go Programming Language website.)

```bash
wget https://storage.googleapis.com/golang/go1.9.3.linux-amd64.tar.gz # Download the Go installer.
sudo tar -C /usr/local -xzf ./go1.9.3.linux-amd64.tar.gz              # Install Go.
rm ./go1.9.3.linux-amd64.tar.gz                                       # Delete the Go installer, as you no longer need it.
```

After you install Go, add the path to the Go binary to your :code:`PATH` environment variable. To do this, add the
following code to the end of your shell profile file (for example, `~/.bashrc` in Amazon Linux, assuming you chose the
option to **Create a new instance for environment (EC2)**, earlier in this topic), and then save the file.

```bash
 PATH=$PATH:/usr/local/go/bin
```

After you save the file, source the `~/.bashrc` file so that the terminal can now find the Go binary you just
referenced. To do this, run the following command. (This command assumes you chose the option to **Create a new instance
for environment (EC2)**, earlier in this topic.)

```bash
. ~/.bashrc
```

## Download Example Code

Use the terminal you opened in the previous step to download example code for the {{% alias sdk-go %}} into the {{%
alias service=AC9 %}} development environment.

To do this, run the following command. This command downloads a copy of all of the code examples used in the official
AWS SDK documentation into your environment's root directory.

```bash
git clone https://github.com/awsdocs/aws-doc-sdk-examples.git
```

To find code examples for the {{% alias sdk-go %}}, use the **Environment** window to open the
`ENVIRONMENT_NAME/aws-doc-sdk-examples/go/example_code` directory, where `ENVIRONMENT_NAME` is the name of your
development environment.

## Run Example Code

To run code in your {{% alias service=AC9 %}} development environment, see
[Run Your Code](https://docs.aws.amazon.com/cloud9/latest/user-guide/build-run-debug.html#build-run-debug-run) in the
{{% alias service=AC9 %}} User Guide.
