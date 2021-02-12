# Go Packages

This is a small collection of packages I created for my golang projects.

## Packages

### Logger (pkg/logger)

The logger package makes it easy to create a logger and use it across multiple files/packages.
It only creates **one instance** and after getting called again returns the pointer to an already created logger.

With this, it makes it really convinient and easy to just import the package and use it without any big concernse. The best thing is taht you an use the logger across all the packages without needing to worry about passing the pointer from package to package.

For more information, head over to the ```docs/logger```, everything from **configuration** to **usage** is there explained in full detail.
