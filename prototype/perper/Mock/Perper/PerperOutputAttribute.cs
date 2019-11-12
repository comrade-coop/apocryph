using System;

namespace Mock.Perper
{
    [System.AttributeUsage(System.AttributeTargets.Parameter | System.AttributeTargets.ReturnValue)]
    public class PerperOutput : Attribute
    {
        public string Name;

        public PerperOutput(string name)
        {
            Name = name;
        }
    }
}
