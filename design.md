# Stubs design

The goal of this library is to be able to generate somewhat good looking random data for structures defined in an openapi or jsonschema document. 
These can be used to provide an api with stubs so that you can verify valid and invalid request/responses without an actually meaningful implementation. 
Some areas where this is useful are collaboration between different teams so that all teams can do work without having to wait on the delivery of the completed implementation of the API.
A second area where this is obviously useful is for doing automated tests and filling up datastructures for use in those tests.


In the openapi 2.0 specification there are 2 families of types that can be used as targets for a random data generator.

1. Path Item related types
   1. Non-Schema Parameters
   2. Response headers
   3. Collection Items on parameters or headers  
2. Schemas

These 2 families warrant slightly different strategies for getting the necessary parameters for a datagenerator.

## Functionality

To generate meaningful random data the library looks at the schema and based on properties for that schema it picks an appropriate random value generator.
This means that when it is generating a string for an object property that represents a first name it should pick a data generator that generates first names.
Similarly should a property be called city it should generate a city name.

It's possible that the validity of a particular property is dependent on the value of another property, for example a creation date is typically never larger than a modified date.

A data generator for a schema should be able to generate valid and invalid data, ideally there is control for specifying the value is invalid for a particular validation. This will help in tests to validate error codes

## How does it work

The main components in this library are a registry of datagenerators so they are addressable by keys, aliases for those keys to aid with inferring which datagenerator to use.
And of course a value generator, which will generate either a valid or an invalid value.

In the openapi API document and in a json schema document there is a vendor extension that can be used to customize the generation process.

### Value generator function

There are 2 types of data enerators but they both have the same function signature.
There are value generators and there are composite value generators.

* A value generator generates a single value for a simple type.
* A composite generator generates a value that is built out of one or more generators eg. object schemas with properties.

The signature for the value generator is:

```go
type ValueGenerator func(GeneratorOpts) (interface{}, error)
```

A value generator is the innermost component in the library and the arguments to the function are used to configure the generato
The generation process for a generator is configured through a GeneratorOpts interface.

### GeneratorOpts interface

The generator options describe the type and potentialy the format for the value that needs to be generated.
In addition to the type information it also captures the field name or definition name of the value that needs to be generated.

### StubMode

The stubmode bitmask allows for configuring which validations should fail for a given value generator.

### Generator

The generator is the main entry point for the librare and it's Generate method is what will generate the random value for the descriptor.
