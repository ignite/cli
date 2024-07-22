# Go Path Environment Setup

This guide provides instructions on how to properly configure your Go PATH environment variables. This is especially useful if you encounter issues like the one described in [GitHub Issue #2291](https://github.com/ignite/cli/issues/2291).

## Prerequisites

Ensure you have Go installed on your system. You can download it from [golang.org](https://golang.org/dl/).

## Steps to Configure Go PATH Environment Variables

1. **Open your shell configuration file**:

   Depending on the shell you are using, open the appropriate configuration file. You might need to use `sudo` to update the file.

   - For `zsh` (Z shell):
     ```sh
     sudo nano ~/.zshrc
     ```

   - For `bash`:
     ```sh
     sudo nano ~/.bashrc
     ```

   - For `profile`:
     ```sh
     sudo nano ~/.profile
     ```

2. **Add the following code snippet**:

   Append the following lines to the end of the opened configuration file to set up your Go PATH environment variables:

   ```sh
   export GOPATH=$HOME/go
   export PATH=$PATH:$GOPATH/bin

   
3. **Apply the changes**:
   
   After saving the changes to your configuration file, you need to source it to apply the changes:

   - For `zsh` (Z shell):
  
        ```sh
        source ~/.zshrc 
        ```
    - For `bash`:
        ```sh
        source ~/.bashrc
        ```
    - For `profile`:
        ```sh
        source ~/.profile
        ```
4. **Verify the configuration**:

    To ensure that the environment variables are set correctly, run the following commands:

    ```sh
        echo $GOPATH
        echo $PATH
    ```



## Troubleshooting

**If you encounter any issues**:

- Double-check that you have correctly edited and sourced the appropriate configuration file.

- Ensure that there are no typos in the environment variable names or paths.

- Restart your terminal session or computer if the changes do not seem to take effect.
