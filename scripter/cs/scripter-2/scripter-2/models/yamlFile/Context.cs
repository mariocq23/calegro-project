using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace scripter.models.yamlFile
{
    public class Context
    {
        public string Name {  get; set; }
        public PerformedAction Action { get; set; }
        public Dictionary<string,string> CustomProperties { get; set; }
    }
}
