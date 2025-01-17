using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace scripter.models.yamlFile
{
    public class ActionConfiguration
    {
        public string Name {  get; set; }
        public string MainFunction { get; set; }
        public string Api {  get; set; }
        public bool DisplayOutputConsole { get; set; }
        public Platform Platform { get; set; }
        public string Path {  get; set; }
        public Dictionary<string,string> EnvironmentVariables { get; set; }

    }
}
