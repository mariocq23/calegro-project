using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace scripter.models.yamlFile
{
    public class Header
    {
        public List<string> Imports {  get; set; }
        public string Id { get; set; }
        public string Name { get; set; }
        public string Inherits { get; set; }
        public List<string> Implements { get; set; }
        
    }
}
