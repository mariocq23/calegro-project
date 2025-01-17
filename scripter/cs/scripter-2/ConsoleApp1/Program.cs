// See https://aka.ms/new-console-template for more information
using scripter;
using scripter_2.dtos;

Console.WriteLine("Hello, World!");

var algo = YamlReader.GetYamlFile<YamlFileDto>("C:/git/calegro-project/examples/vikings/vikings.yaml");

Console.WriteLine(algo);