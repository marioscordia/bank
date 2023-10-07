# Banking

This is a banking REST API that allows user to sign up, sign in, transact money and ask for a loan. It also implements administration, where admin acts as bank which decides whether to accept or reject a loan request.

## Note

- There is only one administrator, which is created at initialization of the database.
  Phone : **87779991100**
  Password : **Cheburek**

## Technologies

Project was built using Fiber as a web framework and PostgreSQL as a database.

## Installation

To run this project, you need to have Docker installed on your machine. Then follow these steps:

1. Pull docker image for PostgreSQL, configure and run it.
2. Navigate to the project directory.
3. Run the following command to build the Docker image:

```bash
make build
```

4. Run the following command to start the Docker container:

```bash
make start
```

5. Open your web browser and go to http://localhost: {portnumber}.

If you want to run without Docker:

```bash
go run .
```
