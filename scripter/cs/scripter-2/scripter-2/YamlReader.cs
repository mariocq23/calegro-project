using scripter.models.yamlFile;
using YamlDotNet.RepresentationModel;
using YamlDotNet.Serialization;
using YamlDotNet.Serialization.NamingConventions;

namespace scripter
{
    public static class YamlReader
    {
        public static T GetYamlFile<T>(string filePath)
        {
            var text = File.ReadAllText(filePath);

            var deserializer = new DeserializerBuilder().Build();
            return deserializer.Deserialize<T>(text);
        }
    }
}
